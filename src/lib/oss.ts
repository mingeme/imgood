import { env } from '$env/dynamic/private';
import { S3Client } from '@aws-sdk/client-s3';

export const client = new S3Client({
  endpoint: `https://${env.DOMAIN}`,
  region: env.REGION,
  credentials: {
    accessKeyId: env.ACCESS_KEY_ID,
    secretAccessKey: env.SECRET_ACCESS_KEY,
  },
});
