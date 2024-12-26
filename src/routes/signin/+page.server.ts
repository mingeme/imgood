import { AuthApiError } from '@supabase/supabase-js';
import { fail, redirect } from '@sveltejs/kit';

export const actions = {
  default: async ({ request, locals: { supabase } }) => {
    const formData = await request.formData();
    const data = Object.fromEntries(formData);
    const { error } = await supabase.auth.signInWithPassword(data);

    if (error) {
      if (error instanceof AuthApiError && error.staartus === 400) {
        return fail(400, {
          error: 'Invalid credentials',
        });
      }
      return fail(500, {
        message: 'Server error. Try again later.',
      });
    }

    redirect(303, '/');
  },
};
