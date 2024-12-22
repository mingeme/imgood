import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';

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
    // 拼接公开 URL，并固定添加查询参数 w=50, h=50, mode=clip
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
