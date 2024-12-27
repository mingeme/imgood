import { env } from '$env/dynamic/private';
import { client } from '$lib/oss';
import { PutObjectCommand } from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { fail } from '@sveltejs/kit';

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
    const name = data.name;
    const timestampBase62 = toBase62(Date.now());
    const key = `${new Date().getFullYear()}/${new Date().getMonth() + 1}/${new Date().getDate()}/${timestampBase62}`;

    const putCmd = new PutObjectCommand({
      Bucket: env.BUCKET,
      ResponseContentDisposition: `attachment; filename=${name}`,
      Key: key,
    });

    const url = await getSignedUrl(client, putCmd, {
      expiresIn: 60 * 5,
      unsignableHeaders: new Set(['content-disposition']),
    });

    const { error } = await supabase.from('image').insert({
      name,
      user_id: user?.id,
      oss_key: key,
      hash: data.hash,
      file_size: data.size,
    });

    if (error) {
      console.error(error);
      return fail(500, { message: 'Failed to save image info' });
    }

    return { url, key };
  },
}; 