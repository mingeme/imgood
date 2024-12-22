import { redirect, type RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
  await supabase.auth.signOut();
  redirect(303, '/signin');
};
