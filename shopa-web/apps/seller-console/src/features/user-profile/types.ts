export type Profile = {
  displayName: string;
  email?: string;
  phone?: string;
  avatarUrl?: string;
  bio?: string;
  updatedAt?: string;
};

export type Address = {
  id: string;
  label: string;
  contact: string;
  phone: string;
  fullAddress: string;
  isDefault: boolean;
  updatedAt?: string;
};
