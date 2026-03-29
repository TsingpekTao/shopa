export type Profile = {
  displayName: string;
  email?: string;
  phone?: string;
  avatarUrl?: string;
  bio?: string;
  updatedAt?: string;
  version?: number;
  addressBookVersion?: number;
};

export type Address = {
  id: string;
  label: string;
  contact: string;
  phone: string;
  fullAddress: string;
  isDefault: boolean;
  updatedAt?: string;
  version?: number;
  addressBookVersion?: number;
};
