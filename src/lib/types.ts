export interface FileWithPreview {
  file: File;
  preview: string;
  hash: string;
}

export interface UploadResponse {
  url: string;
  // 其他响应字段...
}
