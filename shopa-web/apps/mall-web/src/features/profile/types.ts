export type ProfileImageRef = {
  assetId: number;
  url: string;
};

export type MallProfile = {
  displayName: string;
  avatar: ProfileImageRef;
  profileVersion: number;
  ext: Record<string, string>;
  updatedAt: string | { seconds?: number | string; nanos?: number | string } | undefined;
};

export type MallAddress = {
  isDefault: boolean;
};

export type MallProfileBundle = {
  profile?: MallProfile;
  addresses: MallAddress[];
};
