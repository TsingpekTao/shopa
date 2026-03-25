export type MediaAsset = {
  id: string;
  name: string;
  thumbnail: string;
  status: UploadStatus;
  progress: number;
  uploadedAt?: string;
};

export type UploadStatus = 'idle' | 'uploading' | 'processing' | 'done' | 'failed';

export type UploadJob = {
  uploadId: string;
  status: UploadStatus;
  progress: number;
};
