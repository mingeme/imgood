import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, locals: { supabase, user } }) => {
  if (!user) {
    throw error(401, 'Unauthorized');
  }

  const { data: image, error: err } = await supabase
    .from('image')
    .select('*')
    .eq('id', params.id)
    .eq('user_id', user.id)
    .single();

  if (err) {
    if (err.code === 'PGRST116') {
      throw error(404, 'Image not found');
    }
    console.error(err);
    throw error(500, 'Error fetching image');
  }

  return {
    image,
  };
};
