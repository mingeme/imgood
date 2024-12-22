import { env } from '$env/dynamic/private';
import { PutObjectCommand, S3Client } from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { json, type RequestHandler } from '@sveltejs/kit';

// Base62字符集
const base62Chars
    = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';

// 将数字转换为base62
function toBase62(num: number): string {
  let encoded = '';
  while (num > 0) {
    encoded = base62Chars[num % 62] + encoded;
    num = Math.floor(num / 62);
  }
  return encoded;
}

const client = new S3Client({
  endpoint: env.ENDPOINT,
  region: env.REGION,
  credentials: {
    accessKeyId: env.ACCESS_KEY_ID,
    secretAccessKey: env.SECRET_ACCESS_KEY,
  },
});

export const POST: RequestHandler = async ({ request, locals: { supabase, user } }) => {
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
    return json({
      message: 'Internal Server Error',
    });
  }

  return json({
    url,
    key,
  });
};
