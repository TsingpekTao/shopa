export type UserAddress = {
  addressId: number;
  addressVersion: number;
  label: string;
  receiverName: string;
  receiverPhone: string;
  countryCode: string;
  provinceCode: string;
  provinceName: string;
  cityCode: string;
  cityName: string;
  districtCode: string;
  districtName: string;
  street: string;
  detail: string;
  postalCode: string;
  isDefault: boolean;
  status: number;
};

export type AddressBook = {
  addresses: UserAddress[];
  addressBookVersion: number;
};

export type AddressFormInput = {
  label: string;
  receiverName: string;
  receiverPhone: string;
  provinceCode: string;
  provinceName: string;
  cityCode: string;
  cityName: string;
  districtCode: string;
  districtName: string;
  street: string;
  detail: string;
  postalCode: string;
  setAsDefault?: boolean;
};
