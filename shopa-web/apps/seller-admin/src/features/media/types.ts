export type UploadStatus = "idle" | "uploading" | "processing" | "done" | "failed";

export type DerivedAsset = {
  assetId: number;
  kind?: string;
  mimeType?: string;
  publicUrl?: string;
  processStatus?: string | number;
};

export type MediaAsset = {
  id: string;
  assetId: number;
  name: string;
  thumbnail: string;
  status: UploadStatus;
  progress: number;
  uploadedAt?: string;
  publicUrl?: string;
  processStatus?: string | number;
  derivedAssets?: DerivedAsset[];
};

export type UploadJob = {
  uploadId: string;
  assetId?: number;
  status: UploadStatus;
  progress: number;
  derivedAssets?: DerivedAsset[];
};
