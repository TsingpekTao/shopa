"use client";

import { useEffect, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Button,
  Card,
  Empty,
  Input,
  InputNumber,
  Modal,
  Skeleton,
  Space,
  Switch,
  Tag,
  Typography,
  notification
} from "antd";
import {
  DeleteOutlined,
  EditOutlined,
  HolderOutlined,
  PlusOutlined,
  ReloadOutlined,
  TagsOutlined
} from "@ant-design/icons";
import {
  createSellerStoreCategory,
  deleteSellerStoreCategory,
  fetchSellerWorkbench,
  listSellerStoreCategories,
  sortSellerStoreCategories,
  updateSellerStoreCategory
} from "@/features/seller-shop/api";
import { SellerStoreCategory, SellerWorkbenchResponse } from "@/features/seller-shop/types";

const { Title, Text, Paragraph } = Typography;

type EditorState = {
  mode: "create" | "edit";
  parentId: number;
  category?: SellerStoreCategory | null;
};

type DragState = {
  categoryId: number;
  parentId: number;
};

function countChildren(categories: SellerStoreCategory[]) {
  return categories.reduce((sum, item) => sum + item.children.length, 0);
}

function countVisible(categories: SellerStoreCategory[]) {
  return categories.reduce(
    (sum, item) => sum + (item.isVisible ? 1 : 0) + item.children.reduce((childSum, child) => childSum + (child.isVisible ? 1 : 0), 0),
    0
  );
}

function countProducts(categories: SellerStoreCategory[]) {
  return categories.reduce((sum, item) => sum + item.productCount, 0);
}

function reorderCategories(list: SellerStoreCategory[], draggedId: number, targetId: number) {
  const sourceIndex = list.findIndex((item) => item.id === draggedId);
  const targetIndex = list.findIndex((item) => item.id === targetId);
  if (sourceIndex < 0 || targetIndex < 0 || sourceIndex === targetIndex) {
    return list;
  }

  const next = [...list];
  const [dragged] = next.splice(sourceIndex, 1);
  next.splice(targetIndex, 0, dragged);
  return next.map((item, index) => ({
    ...item,
    sortOrder: (index + 1) * 10
  }));
}

function applySiblingReorder(
  categories: SellerStoreCategory[],
  parentId: number,
  draggedId: number,
  targetId: number
) {
  if (parentId === 0) {
    return reorderCategories(categories, draggedId, targetId);
  }

  return categories.map((category) => {
    if (category.id !== parentId) {
      return category;
    }
    return {
      ...category,
      children: reorderCategories(category.children, draggedId, targetId)
    };
  });
}

function extractSiblings(categories: SellerStoreCategory[], parentId: number) {
  if (parentId === 0) {
    return categories;
  }
  return categories.find((item) => item.id === parentId)?.children ?? [];
}

export default function StoreCategoriesPage() {
  const [editorState, setEditorState] = useState<EditorState | null>(null);
  const [name, setName] = useState("");
  const [sortOrder, setSortOrder] = useState<number>(0);
  const [isVisible, setIsVisible] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [dragState, setDragState] = useState<DragState | null>(null);
  const [localCategories, setLocalCategories] = useState<SellerStoreCategory[] | null>(null);
  const [sorting, setSorting] = useState(false);

  const workbenchQuery = useQuery<SellerWorkbenchResponse>({
    queryKey: ["seller-workbench", "store-categories"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30_000
  });

  const shop = workbenchQuery.data?.shops?.[0];

  const categoriesQuery = useQuery({
    queryKey: ["seller-store-categories", shop?.shopNo],
    queryFn: () => listSellerStoreCategories(shop?.shopNo ?? ""),
    enabled: Boolean(shop?.shopNo),
    staleTime: 15_000
  });

  useEffect(() => {
    setLocalCategories(categoriesQuery.data ?? null);
  }, [categoriesQuery.data]);

  const categories = localCategories ?? categoriesQuery.data ?? [];
  const totalSecondary = useMemo(() => countChildren(categories), [categories]);
  const totalVisible = useMemo(() => countVisible(categories), [categories]);
  const totalProducts = useMemo(() => countProducts(categories), [categories]);

  const openCreate = (parentId = 0) => {
    setEditorState({ mode: "create", parentId });
    setName("");
    setSortOrder(0);
    setIsVisible(true);
  };

  const openEdit = (category: SellerStoreCategory) => {
    setEditorState({ mode: "edit", parentId: category.parentId, category });
    setName(category.name);
    setSortOrder(category.sortOrder);
    setIsVisible(category.isVisible);
  };

  const closeEditor = () => {
    setEditorState(null);
    setName("");
    setSortOrder(0);
    setIsVisible(true);
  };

  const handleSubmit = async () => {
    if (!shop?.shopNo) {
      notification.warning({ message: "当前没有可管理的店铺" });
      return;
    }
    if (!name.trim()) {
      notification.warning({ message: "请填写分类名称" });
      return;
    }

    setSubmitting(true);
    try {
      if (editorState?.mode === "edit" && editorState.category) {
        await updateSellerStoreCategory(shop.shopNo, editorState.category.id, {
          name: name.trim(),
          sortOrder,
          isVisible
        });
        notification.success({ message: "分类已更新" });
      } else {
        await createSellerStoreCategory(shop.shopNo, {
          parentId: editorState?.parentId ?? 0,
          name: name.trim(),
          sortOrder,
          isVisible
        });
        notification.success({ message: editorState?.parentId ? "二级分类已创建" : "一级分类已创建" });
      }
      closeEditor();
      await categoriesQuery.refetch();
    } catch (error) {
      notification.error({
        message: editorState?.mode === "edit" ? "分类更新失败" : "分类创建失败",
        description: error instanceof Error ? error.message : "请稍后重试"
      });
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = (category: SellerStoreCategory) => {
    if (!shop?.shopNo) {
      return;
    }
    Modal.confirm({
      title: "删除店内分类",
      content:
        category.children.length > 0
          ? "该分类下还有二级分类，当前不能删除。请先清空子分类后再操作。"
          : "只有空分类允许删除；如果仍有商品绑定，系统会阻止删除。",
      okButtonProps: { danger: true, disabled: category.children.length > 0 },
      onOk: async () => {
        try {
          await deleteSellerStoreCategory(shop.shopNo, category.id);
          notification.success({ message: "分类已删除" });
          await categoriesQuery.refetch();
        } catch (error) {
          notification.error({
            message: "删除失败",
            description: error instanceof Error ? error.message : "请先解绑或迁移商品后再试"
          });
        }
      }
    });
  };

  const commitReorder = async (parentId: number, draggedId: number, targetId: number) => {
    if (!shop?.shopNo || draggedId === targetId) {
      return;
    }

    const nextCategories = applySiblingReorder(categories, parentId, draggedId, targetId);
    const siblings = extractSiblings(nextCategories, parentId);
    setLocalCategories(nextCategories);
    setSorting(true);

    try {
      await sortSellerStoreCategories(
        shop.shopNo,
        siblings.map((item) => ({
          categoryId: item.id,
          sortOrder: item.sortOrder
        }))
      );
      notification.success({ message: "排序已更新" });
      await categoriesQuery.refetch();
    } catch (error) {
      setLocalCategories(categoriesQuery.data ?? []);
      notification.error({
        message: "排序更新失败",
        description: error instanceof Error ? error.message : "请稍后重试"
      });
    } finally {
      setSorting(false);
      setDragState(null);
    }
  };

  const handleDragStart = (categoryId: number, parentId: number) => {
    setDragState({ categoryId, parentId });
  };

  const handleDrop = async (targetId: number, parentId: number) => {
    if (!dragState || dragState.parentId !== parentId) {
      setDragState(null);
      return;
    }
    await commitReorder(parentId, dragState.categoryId, targetId);
  };

  if (workbenchQuery.isLoading) {
    return <Skeleton active paragraph={{ rows: 10 }} />;
  }

  if (!shop) {
    return (
      <Card bordered={false} style={{ borderRadius: 24 }}>
        <Empty description="当前账号还没有可管理的店铺" />
      </Card>
    );
  }

  return (
    <section className="seller-page" style={{ gap: 20 }}>
      <Card
        bordered={false}
        style={{
          borderRadius: 28,
          border: "1px solid #ead8c5",
          background: "linear-gradient(135deg, #fff8f1 0%, #fffdf9 100%)",
          boxShadow: "0 18px 42px rgba(103, 66, 26, 0.08)"
        }}
      >
        <Space direction="vertical" size={16} style={{ width: "100%" }}>
          <div style={{ display: "flex", justifyContent: "space-between", gap: 16, flexWrap: "wrap" }}>
            <div style={{ display: "grid", gap: 8 }}>
              <Text style={{ fontSize: 12, letterSpacing: 2, textTransform: "uppercase", color: "#b56a2b" }}>
                In-Shop Taxonomy
              </Text>
              <Title level={3} style={{ margin: 0 }}>
                店内分类管理
              </Title>
              <Paragraph style={{ margin: 0, color: "#7e6652", maxWidth: 760 }}>
                这里管理店铺自己的一级/二级分类。买家店铺页左侧导航、顶部二级 tabs，以及商品发布页的店内分类选择器，都会以这里的真实分类为准。
              </Paragraph>
            </div>

            <Space wrap>
              <Tag color="gold">{shop.shopDisplayName || shop.shopName}</Tag>
              <Button icon={<ReloadOutlined />} onClick={() => void categoriesQuery.refetch()} loading={categoriesQuery.isFetching && !sorting}>
                刷新
              </Button>
              <Button type="primary" icon={<PlusOutlined />} onClick={() => openCreate(0)} style={{ borderRadius: 12 }}>
                新建一级分类
              </Button>
            </Space>
          </div>

          <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))" }}>
            <Card size="small" style={{ borderRadius: 18, background: "#fffdfb" }}>
              <Text type="secondary">一级分类</Text>
              <div style={{ marginTop: 8, fontSize: 24, fontWeight: 700 }}>{categories.length}</div>
            </Card>
            <Card size="small" style={{ borderRadius: 18, background: "#fffdfb" }}>
              <Text type="secondary">二级分类</Text>
              <div style={{ marginTop: 8, fontSize: 24, fontWeight: 700 }}>{totalSecondary}</div>
            </Card>
            <Card size="small" style={{ borderRadius: 18, background: "#fffdfb" }}>
              <Text type="secondary">可见分类</Text>
              <div style={{ marginTop: 8, fontSize: 24, fontWeight: 700 }}>{totalVisible}</div>
            </Card>
            <Card size="small" style={{ borderRadius: 18, background: "#fffdfb" }}>
              <Text type="secondary">聚合商品数</Text>
              <div style={{ marginTop: 8, fontSize: 24, fontWeight: 700 }}>{totalProducts}</div>
            </Card>
          </div>
        </Space>
      </Card>

      <Card bordered={false} style={{ borderRadius: 24 }}>
        <div style={{ display: "flex", justifyContent: "space-between", gap: 12, flexWrap: "wrap", marginBottom: 18 }}>
          <div>
            <Title level={4} style={{ marginBottom: 4 }}>
              分类树
            </Title>
            <Text type="secondary">
              拖拽同层级分类即可重新排序。一级分类的商品数会自动包含其下所有二级分类；数据库里仍只保存各分类自己的直绑数。
            </Text>
          </div>
          {sorting ? <Tag color="processing">正在保存排序</Tag> : <Tag>拖拽排序已开启</Tag>}
        </div>

        {categoriesQuery.isLoading ? <Skeleton active paragraph={{ rows: 8 }} /> : null}

        {!categoriesQuery.isLoading && categories.length === 0 ? (
          <Empty description="店铺还没有配置店内分类" image={Empty.PRESENTED_IMAGE_SIMPLE}>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => openCreate(0)}>
              先创建一级分类
            </Button>
          </Empty>
        ) : null}

        <div style={{ display: "grid", gap: 16 }}>
          {categories.map((category) => (
            <Card
              key={category.id}
              size="small"
              draggable
              onDragStart={() => handleDragStart(category.id, 0)}
              onDragOver={(event) => event.preventDefault()}
              onDrop={() => void handleDrop(category.id, 0)}
              style={{
                borderRadius: 22,
                border: dragState?.categoryId === category.id ? "1px solid #ffb17e" : "1px solid #efdfce",
                background: "#fffaf6",
                cursor: "grab"
              }}
              title={
                <div style={{ display: "flex", alignItems: "center", gap: 10, flexWrap: "wrap" }}>
                  <HolderOutlined style={{ color: "#b4672f" }} />
                  <TagsOutlined />
                  <Text strong>{category.name}</Text>
                  <Tag color={category.isVisible ? "green" : "default"}>{category.isVisible ? "前台展示" : "已隐藏"}</Tag>
                  <Tag color="blue">商品数 {category.productCount}</Tag>
                  <Tag>排序 {category.sortOrder}</Tag>
                </div>
              }
              extra={
                <Space wrap size={8}>
                  <Button size="small" icon={<PlusOutlined />} onClick={() => openCreate(category.id)}>
                    新建二级
                  </Button>
                  <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(category)}>
                    编辑
                  </Button>
                  <Button
                    size="small"
                    danger
                    icon={<DeleteOutlined />}
                    disabled={category.children.length > 0}
                    onClick={() => handleDelete(category)}
                  >
                    删除
                  </Button>
                </Space>
              }
            >
              {category.children.length > 0 ? (
                <div style={{ display: "grid", gap: 12 }}>
                  {category.children.map((child) => (
                    <div
                      key={child.id}
                      draggable
                      onDragStart={() => handleDragStart(child.id, category.id)}
                      onDragOver={(event) => event.preventDefault()}
                      onDrop={() => void handleDrop(child.id, category.id)}
                      style={{
                        borderRadius: 18,
                        border: dragState?.categoryId === child.id ? "1px solid #ffb17e" : "1px solid #f1e5d9",
                        background: "#fff",
                        padding: 14,
                        display: "flex",
                        justifyContent: "space-between",
                        gap: 16,
                        flexWrap: "wrap",
                        cursor: "grab"
                      }}
                    >
                      <div style={{ display: "grid", gap: 6 }}>
                        <Space wrap>
                          <HolderOutlined style={{ color: "#b4672f" }} />
                          <Text strong>{child.name}</Text>
                          <Tag color={child.isVisible ? "green" : "default"}>{child.isVisible ? "前台展示" : "已隐藏"}</Tag>
                        </Space>
                        <Text type="secondary">直绑商品数 {child.productCount}，排序 {child.sortOrder}</Text>
                      </div>

                      <Space wrap>
                        <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(child)}>
                          编辑
                        </Button>
                        <Button size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(child)}>
                          删除
                        </Button>
                      </Space>
                    </div>
                  ))}
                </div>
              ) : (
                <div
                  style={{
                    borderRadius: 18,
                    border: "1px dashed #e7d5c4",
                    background: "#fff",
                    padding: 18,
                    color: "#7e6652"
                  }}
                >
                  当前没有二级分类。这个一级分类可直接作为叶子分类绑定商品，也可以继续新增二级分类。
                </div>
              )}
            </Card>
          ))}
        </div>
      </Card>

      <Modal
        open={Boolean(editorState)}
        title={editorState?.mode === "edit" ? "编辑店内分类" : editorState?.parentId ? "新建二级分类" : "新建一级分类"}
        onCancel={closeEditor}
        onOk={() => void handleSubmit()}
        okText={editorState?.mode === "edit" ? "保存修改" : "创建分类"}
        confirmLoading={submitting}
      >
        <Space direction="vertical" size={18} style={{ width: "100%" }}>
          <label style={{ display: "grid", gap: 8 }}>
            <Text strong>分类名称</Text>
            <Input value={name} onChange={(event) => setName(event.target.value)} placeholder="例如：上衣 / 下装 / 外套 / 鞋子" />
          </label>

          <label style={{ display: "grid", gap: 8 }}>
            <Text strong>排序值</Text>
            <InputNumber value={sortOrder} onChange={(value) => setSortOrder(typeof value === "number" ? value : 0)} style={{ width: "100%" }} />
          </label>

          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              gap: 12,
              borderRadius: 16,
              padding: "12px 14px",
              background: "#fff8f2"
            }}
          >
            <div style={{ display: "grid", gap: 4 }}>
              <Text strong>前台可见</Text>
              <Text type="secondary">关闭后分类不会出现在买家店铺页，但后台仍保留。</Text>
            </div>
            <Switch checked={isVisible} onChange={setIsVisible} />
          </div>
        </Space>
      </Modal>
    </section>
  );
}
