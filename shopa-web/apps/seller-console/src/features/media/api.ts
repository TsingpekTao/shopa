import { MediaAsset, UploadStatus, UploadJob } from './types';

const delay = (ms: number) => new Promise<void>((resolve) => setTimeout(resolve, ms));

export async function fetchMediaLibrary(): Promise<MediaAsset[]> {
  await delay(400);
  return [
    {
      id: 'asset-1',
      name: 'Brand Authorization.pdf',
      thumbnail: '/images/media/preview-brand.jpg',
      status: 'processing',
      progress: 58,
      uploadedAt: '2026-03-24T10:33:00Z',
    },
    {
      id: 'asset-2',
      name: 'Storefront Cover.png',
      thumbnail: '/images/media/storefront.png',
      status: 'done',
      progress: 100,
      uploadedAt: '2026-03-24T09:12:00Z',
    },
  ];
}

export async function startUpload(): Promise<UploadJob> {
  await delay(250);
  return { uploadId: `upload-${Date.now()}`, status: 'uploading', progress: 12 };
}

export async function pollUploadStatus(uploadId: string): Promise<UploadJob> {
  await delay(350);
  const tick = Math.min(100, Math.floor((Date.now() / 1000) % 100));
  const status: UploadStatus = tick < 60 ? 'uploading' : tick < 90 ? 'processing' : 'done';
  return { uploadId, status, progress: tick };
}

// TODO: replace mocks with real upload endpoint: POST /v1/media/upload + polling /v1/media/{id}/status
