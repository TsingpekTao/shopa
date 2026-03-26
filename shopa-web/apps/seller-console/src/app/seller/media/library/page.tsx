import { Metadata } from 'next';
import MediaSection from '@/features/media/MediaSection';

export const metadata: Metadata = {
  title: 'Seller Media Library',
  description: 'Upload and manage seller assets before publishing a product.',
};

const Page = () => (
  <main className="media-layout">
    <MediaSection />
  </main>
);

export default Page;
