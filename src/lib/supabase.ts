import type { Database } from './database.types';
import { env } from '$env/dynamic/public';
import { createClient } from '@supabase/supabase-js';

const supabase = createClient<Database>(
  env.PUBLIC_SUPABASE_URL,
  env.PUBLIC_SUPABASE_ANON_KEY,
);

export default supabase;
