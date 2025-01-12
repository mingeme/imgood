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

  return {
    images: data,
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
      console.error('Delete image from OSS failed', error);
      return fail(500, { message: 'Delete image from OSS failed' });
    }

    return { success: true };
  },
};
