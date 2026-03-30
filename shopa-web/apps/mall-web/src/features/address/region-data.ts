import { areaList } from "@vant/area-data";

export type RegionDistrict = {
  code: string;
  name: string;
};

export type RegionCity = {
  code: string;
  name: string;
  districts: RegionDistrict[];
};

export type RegionProvince = {
  code: string;
  name: string;
  cities: RegionCity[];
};

type RegionOption = {
  code: string;
  name: string;
};

const DIRECT_CITY_PROVINCES = new Set(["11", "12", "31", "50"]);

const provinceNameMap = toSortedMap(areaList.province_list);
const cityNameMap = toSortedMap(areaList.city_list);
const districtNameMap = toSortedMap(areaList.county_list);

const provinces: RegionProvince[] = buildProvinces();
const provinceMap = new Map<string, RegionProvince>(provinces.map((item) => [item.code, item]));
const cityMap = new Map<string, RegionCity>();
const districtMap = new Map<string, RegionDistrict>();

for (const province of provinces) {
  for (const city of province.cities) {
    cityMap.set(city.code, city);
    for (const district of city.districts) {
      districtMap.set(district.code, district);
    }
  }
}

function toSortedMap(input: Record<string, string>): Map<string, string> {
  const entries = Object.entries(input)
    .map(([code, name]) => [String(code), String(name)] as const)
    .sort((a, b) => Number(a[0]) - Number(b[0]));
  return new Map(entries);
}

function buildProvinces(): RegionProvince[] {
  const result: RegionProvince[] = [];

  for (const [provinceCode, provinceName] of provinceNameMap.entries()) {
    const provincePrefix = provinceCode.slice(0, 2);
    const provinceCities: RegionCity[] = [];

    for (const [cityCode, cityName] of cityNameMap.entries()) {
      if (cityCode.slice(0, 2) !== provincePrefix) {
        continue;
      }
      const districts = collectDistricts(provinceCode, cityCode);
      provinceCities.push({ code: cityCode, name: cityName, districts });
    }

    if (provinceCities.length === 0) {
      const fallbackCityCode = `${provincePrefix}0100`;
      const fallbackCityName = provinceName;
      provinceCities.push({
        code: fallbackCityCode,
        name: fallbackCityName,
        districts: collectDistricts(provinceCode, fallbackCityCode)
      });
    }

    result.push({
      code: provinceCode,
      name: provinceName,
      cities: provinceCities
    });
  }

  return result;
}

function collectDistricts(provinceCode: string, cityCode: string): RegionDistrict[] {
  const result: RegionDistrict[] = [];
  const provincePrefix = provinceCode.slice(0, 2);
  const cityPrefix = cityCode.slice(0, 4);

  for (const [districtCode, districtName] of districtNameMap.entries()) {
    if (districtCode.slice(0, 2) !== provincePrefix) {
      continue;
    }

    if (DIRECT_CITY_PROVINCES.has(provincePrefix)) {
      if (districtCode.slice(0, 2) !== cityCode.slice(0, 2)) {
        continue;
      }
    } else if (districtCode.slice(0, 4) !== cityPrefix) {
      continue;
    }

    result.push({ code: districtCode, name: districtName });
  }

  if (result.length === 0) {
    const cityName = cityNameMap.get(cityCode) ?? cityCode;
    result.push({
      code: `${cityCode.slice(0, 4)}01`,
      name: cityName
    });
  }

  return result;
}

export function getProvinceOptions(): RegionOption[] {
  return provinces.map((item) => ({ code: item.code, name: item.name }));
}

export function getCityOptions(provinceCode: string): RegionOption[] {
  const province = provinceMap.get(provinceCode);
  if (!province) {
    return [];
  }
  return province.cities.map((item) => ({ code: item.code, name: item.name }));
}

export function getDistrictOptions(provinceCode: string, cityCode: string): RegionOption[] {
  const province = provinceMap.get(provinceCode);
  if (!province) {
    return [];
  }
  const city = province.cities.find((item) => item.code === cityCode);
  if (!city) {
    return [];
  }
  return city.districts.map((item) => ({ code: item.code, name: item.name }));
}

export function findProvinceByCode(code: string): RegionProvince | undefined {
  return provinceMap.get(code);
}

export function findCityByCode(provinceCode: string, cityCode: string): RegionCity | undefined {
  const province = provinceMap.get(provinceCode);
  if (!province) {
    return undefined;
  }
  return province.cities.find((item) => item.code === cityCode);
}

export function findDistrictByCode(provinceCode: string, cityCode: string, districtCode: string): RegionDistrict | undefined {
  const city = findCityByCode(provinceCode, cityCode);
  if (!city) {
    return undefined;
  }
  return city.districts.find((item) => item.code === districtCode);
}

export function findProvinceCodeByName(name: string): string {
  const normalized = name.trim();
  if (!normalized) {
    return "";
  }
  for (const [code, text] of provinceNameMap.entries()) {
    if (text === normalized) {
      return code;
    }
  }
  return "";
}

export function findCityCodeByName(provinceCode: string, name: string): string {
  const normalized = name.trim();
  if (!normalized) {
    return "";
  }
  const province = provinceMap.get(provinceCode);
  if (!province) {
    return "";
  }
  return province.cities.find((item) => item.name === normalized)?.code ?? "";
}

export function findDistrictCodeByName(provinceCode: string, cityCode: string, name: string): string {
  const normalized = name.trim();
  if (!normalized) {
    return "";
  }
  const city = findCityByCode(provinceCode, cityCode);
  if (!city) {
    return "";
  }
  return city.districts.find((item) => item.name === normalized)?.code ?? "";
}

