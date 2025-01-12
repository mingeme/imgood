export interface FileWithPreview {
  file: File;
  preview: string;
  hash: string;
  status: 'uploading' | 'success' | 'error' | 'duplicated' | 'pending';
}

export interface UploadResponse {
  url: string;
}
