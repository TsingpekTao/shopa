export type MallOverviewRole = {
  roleCode: string;
  scopeTypeCode: string;
  scopeNo: string;
};

export type MallOverview = {
  userId: number;
  accountStatusCode: string;
  roles: MallOverviewRole[];
  displayName: string;
  avatarUrl: string;
  points: number;
  partial: boolean;
  degradedFields: string[];
};
