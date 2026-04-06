export type BuyerCategoryDescriptor = {
  key: string;
  labelZh: string;
  labelEn: string;
  groupKey: string;
  groupLabelZh: string;
  groupLabelEn: string;
};

const KNOWN_CATEGORY_BY_ID: Record<number, BuyerCategoryDescriptor> = {
  10001: {
    key: "denim",
    labelZh: "\u725b\u4ed4\u88e4",
    labelEn: "Denim",
    groupKey: "bottoms",
    groupLabelZh: "\u88e4\u88c5",
    groupLabelEn: "Bottoms"
  },
  10002: {
    key: "leather-jackets",
    labelZh: "\u76ae\u5939\u514b",
    labelEn: "Leather jackets",
    groupKey: "outerwear",
    groupLabelZh: "\u5916\u5957",
    groupLabelEn: "Outerwear"
  },
  10003: {
    key: "wide-leg-pants",
    labelZh: "\u9614\u817f\u88e4",
    labelEn: "Wide-leg pants",
    groupKey: "bottoms",
    groupLabelZh: "\u88e4\u88c5",
    groupLabelEn: "Bottoms"
  },
  10004: {
    key: "utility-jackets",
    labelZh: "\u5de5\u88c5\u5916\u5957",
    labelEn: "Utility jackets",
    groupKey: "outerwear",
    groupLabelZh: "\u5916\u5957",
    groupLabelEn: "Outerwear"
  },
  100001: {
    key: "running-shoes",
    labelZh: "\u8dd1\u6b65\u978b",
    labelEn: "Running shoes",
    groupKey: "shoes",
    groupLabelZh: "\u978b\u5b50",
    groupLabelEn: "Shoes"
  }
};

const KEYWORD_RULES: Array<{
  test: RegExp;
  descriptor: BuyerCategoryDescriptor;
}> = [
  {
    test: /(\u978b|sneaker|shoe|runner|running)/i,
    descriptor: {
      key: "shoes",
      labelZh: "\u978b\u5b50",
      labelEn: "Shoes",
      groupKey: "shoes",
      groupLabelZh: "\u978b\u5b50",
      groupLabelEn: "Shoes"
    }
  },
  {
    test: /(\u5305|\u80cc\u5305|\u624b\u63d0|\u659c\u633a|bag|backpack|tote)/i,
    descriptor: {
      key: "bags",
      labelZh: "\u5305\u888b",
      labelEn: "Bags",
      groupKey: "bags",
      groupLabelZh: "\u5305\u888b",
      groupLabelEn: "Bags"
    }
  },
  {
    test: /(\u88e4|\u77ed\u88e4|\u725b\u4ed4|\u88d9|pants|jeans|shorts|skirt)/i,
    descriptor: {
      key: "bottoms",
      labelZh: "\u88e4\u88c5",
      labelEn: "Bottoms",
      groupKey: "bottoms",
      groupLabelZh: "\u88e4\u88c5",
      groupLabelEn: "Bottoms"
    }
  },
  {
    test: /(\u5916\u5957|\u5939\u514b|\u98ce\u8863|\u5927\u8863|\u76ae\u8863|coat|jacket|outerwear|blazer)/i,
    descriptor: {
      key: "outerwear",
      labelZh: "\u5916\u5957",
      labelEn: "Outerwear",
      groupKey: "outerwear",
      groupLabelZh: "\u5916\u5957",
      groupLabelEn: "Outerwear"
    }
  },
  {
    test: /(\u536b\u8863|t\u6064|\u886c\u886b|\u9488\u7ec7|\u4e0a\u8863|shirt|tee|hoodie|sweater|top)/i,
    descriptor: {
      key: "tops",
      labelZh: "\u4e0a\u8863",
      labelEn: "Tops",
      groupKey: "tops",
      groupLabelZh: "\u4e0a\u8863",
      groupLabelEn: "Tops"
    }
  },
  {
    test: /(\u5e3d|\u889c|\u8170\u5e26|\u56f4\u5dfe|\u914d\u4ef6|hat|sock|belt|scarf|accessor)/i,
    descriptor: {
      key: "accessories",
      labelZh: "\u914d\u4ef6",
      labelEn: "Accessories",
      groupKey: "accessories",
      groupLabelZh: "\u914d\u4ef6",
      groupLabelEn: "Accessories"
    }
  }
];

const FALLBACK_DESCRIPTOR: BuyerCategoryDescriptor = {
  key: "featured",
  labelZh: "\u7cbe\u9009\u5355\u54c1",
  labelEn: "Featured",
  groupKey: "featured",
  groupLabelZh: "\u7cbe\u9009",
  groupLabelEn: "Featured"
};

function normalizeText(input: string) {
  return input.trim().toLowerCase();
}

function slugify(input: string) {
  const normalized = normalizeText(input)
    .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-")
    .replace(/^-+|-+$/g, "");
  return normalized || FALLBACK_DESCRIPTOR.key;
}

export function resolveBuyerCategoryDescriptor(input: {
  categoryId?: number;
  categoryName?: string;
  title?: string;
}): BuyerCategoryDescriptor {
  if (input.categoryId && input.categoryId in KNOWN_CATEGORY_BY_ID) {
    return KNOWN_CATEGORY_BY_ID[input.categoryId];
  }

  const explicitName = input.categoryName?.trim();
  if (explicitName) {
    const byKeyword = KEYWORD_RULES.find((rule) => rule.test.test(explicitName))?.descriptor;
    if (byKeyword) {
      return {
        ...byKeyword,
        key: slugify(explicitName),
        labelZh: explicitName,
        labelEn: explicitName
      };
    }

    return {
      key: slugify(explicitName),
      labelZh: explicitName,
      labelEn: explicitName,
      groupKey: slugify(explicitName),
      groupLabelZh: explicitName,
      groupLabelEn: explicitName
    };
  }

  const title = input.title?.trim() || "";
  const byTitle = KEYWORD_RULES.find((rule) => rule.test.test(title))?.descriptor;
  return byTitle ?? FALLBACK_DESCRIPTOR;
}
