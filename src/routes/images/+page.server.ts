import { env } from '$env/dynamic/private';
import { client } from '$lib/oss';
import { DeleteObjectCommand } from '@aws-sdk/client-s3';
import { fail, redirect } from '@sveltejs/kit';

export async function load({ locals: { supabase, user } }) {
  const { data, error } = await supabase.from('image')
    .select()
    .eq('user_id', user?.id)
    .order('created_at', { ascending: false });

  if (error) {
    console.error(error);
    redirect(303, '/error');
  }

  const imagesWithUrls = await Promise.all(data.map(async (image) => {
    const publicUrl = `https://${env.BUCKET}.${env.DOMAIN}/${image.oss_key}`;
    const previewUrl = `${publicUrl}?w=50&h=50&mode=clip`;
    return {
      ...image,
      publicUrl,
      previewUrl,
    };
  }));

  return {
    images: imagesWithUrls,
  };
}

export const actions = {
  delete: async ({ request, locals: { supabase, user } }) => {
    const formData = await request.formData();
    const ossKey = formData.get('ossKey') as string;

    const { error: supabaseError } = await supabase
      .from('image')
      .delete()
      .eq('user_id', user?.id)
      .eq('oss_key', ossKey);
    if (supabaseError) {
      return fail(500, { message: 'Failed to delete image from Supabase' });
    }

    try {
      await client.send(new DeleteObjectCommand({
        Bucket: env.BUCKET,
        Key: ossKey,
      }));
    }
    catch (error) {
      console.error('删除 OSS 对象失败:', error);
      return fail(500, { message: '从 OSS 删除图片失败' });
    }

    return { success: true };
  },
};
