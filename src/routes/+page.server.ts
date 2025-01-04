import { env } from '$env/dynamic/private';
import { client } from '$lib/oss';
import { PutObjectCommand } from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { error, fail } from '@sveltejs/kit';

const base62Chars = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';

function toBase62(num: number): string {
  let encoded = '';
  while (num > 0) {
    encoded = base62Chars[num % 62] + encoded;
    num = Math.floor(num / 62);
  }
  return encoded;
}

export const actions = {
  upload: async ({ request, locals: { supabase, user } }) => {
    const data = await request.json();
    const { name, hash } = data;
    const timestampBase62 = toBase62(Date.now());
    const key = `${new Date().getFullYear()}/${new Date().getMonth() + 1}/${new Date().getDate()}/${timestampBase62}`;

    const { data: existingImages, error: checkError } = await supabase
      .from('image')
      .select('id')
      .eq('hash', hash)
      .single();

    if (checkError && checkError.code !== 'PGRST116') {
      console.error(checkError);
      return fail(500, { message: 'Failed to check for existing image' });
    }

    if (existingImages) {
      return fail(400, { message: 'This image has already been uploaded.' });
    }

    const putCmd = new PutObjectCommand({
      Bucket: env.BUCKET,
      ResponseContentDisposition: `attachment; filename=${name}`,
      Key: key,
    });

    const url = await getSignedUrl(client, putCmd, {
      expiresIn: 60 * 5,
      unsignableHeaders: new Set(['content-disposition']),
    });

    const { error: checkErr } = await supabase.from('image').insert({
      name,
      user_id: user?.id,
      oss_key: key,
      hash,
      file_size: data.size,
    });

    if (checkErr) {
      console.error(checkErr);
      return error(500, { message: 'Failed to save image info' });
    }

    return { url, key };
  },
};
