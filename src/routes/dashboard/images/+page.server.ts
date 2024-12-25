import { env } from '$env/dynamic/private';
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
    const publicUrl = `https://${env.BUCKET}.s3.bitiful.net/${image.oss_key}?w=50&h=50&mode=clip`;
    return {
      ...image,
      url: publicUrl,
    };
  }));

  return {
    images: imagesWithUrls,
  };
}

export const actions = {
  delete: async ({ request, locals: { supabase, user } }) => {
    const formData = await request.formData();
    const ossKey = formData.get('oss_key') as string;

    const { error: supabaseError } = await supabase
      .from('image')
      .delete()
      .eq('user_id', user?.id)
      .eq('oss_key', ossKey);
    if (supabaseError) {
      return fail(500, { message: 'Failed to delete image from Supabase' });
    }

    // TODO delete oss
    return { success: true };
  },
};
