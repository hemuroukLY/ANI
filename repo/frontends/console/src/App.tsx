import { type ReactNode, useMemo, useState } from 'react'
import { Button, Card, Space, Table, Tag } from 'tdesign-react'
import {
  DashboardIcon,
  ChartBarIcon,
  SettingIcon,
  ServerIcon,
  SearchIcon,
} from 'tdesign-icons-react'
import { InstanceDemoPage } from './demo/InstanceDemoPage'
import { DemoConsolePage } from './demo/DemoConsolePage'

const brandLogo = '/assets/brand/d-logo.png'

type ConsoleView = 'demo-instances' | 'dashboard' | 'models' | 'kb' | 'usage' | 'settings'

const viewTitles: Record<ConsoleView, string> = {
  'demo-instances': '实例工作台',
  dashboard: '运营总览',
  models: '模型管理',
  kb: '知识库',
  usage: '用量报表',
  settings: '平台设置',
}

const navItems: Array<{ value: ConsoleView; label: string; icon: ReactNode }> = [
  { value: 'dashboard', label: '运营总览', icon: <DashboardIcon /> },
  { value: 'demo-instances', label: '实例工作台', icon: <ServerIcon /> },
  { value: 'models', label: '模型管理', icon: <ServerIcon /> },
  { value: 'kb', label: '知识库', icon: <ServerIcon /> },
  { value: 'usage', label: '用量报表', icon: <ChartBarIcon /> },
  { value: 'settings', label: '平台设置', icon: <SettingIcon /> },
]

export function App() {
  const [view, setView] = useState<ConsoleView>('demo-instances')

  if (window.location.pathname === '/demo-console') {
    return <DemoConsolePage />
  }

  return (
    <div className="app-layout">
      <header className="app-header">
        <div className="brand-lockup">
          <img src={brandLogo} alt="常青云" />
          <div>
            <strong>常青云 ANI</strong>
            <span>AI Native Infrastructure</span>
          </div>
        </div>
        <div className="global-search">
          <SearchIcon />
          <span>搜索服务、实例、模型、任务</span>
        </div>
        <div className="topbar-actions">
          <button>支持</button>
          <button>通知</button>
          <button>cn-north-demo</button>
          <button>demo-admin</button>
        </div>
      </header>
      <div className="app-body">
        <aside className="app-aside">
          <div className="service-title">服务导航</div>
          <nav className="service-nav">
            {navItems.map((item) => (
              <button
                key={item.value}
                className={view === item.value ? 'service-nav-active' : ''}
                onClick={() => setView(item.value)}
              >
                {item.icon}
                <span>{item.label}</span>
              </button>
            ))}
          </nav>
          <div className="aside-status">
            <span>Demo Profile</span>
            <strong>M1 Instance Fabric</strong>
          </div>
        </aside>
        <main className="app-content">
          <div className="content-breadcrumb">ANI Console / {viewTitles[view]}</div>
          {view === 'demo-instances' ? <InstanceDemoPage /> : <MockConsolePage view={view} />}
        </main>
      </div>
    </div>
  )
}

function MockConsolePage({ view }: { view: ConsoleView }) {
  const [selected, setSelected] = useState(0)
  const [message, setMessage] = useState('')
  const data = useMemo(() => mockRows(view), [view])
  const row = data[selected] ?? data[0]
  const columns = [
    {
      title: '名称',
      colKey: 'name',
      cell: ({ row: item }: { row: MockRow }) => (
        <button className="link-button" onClick={() => setSelected(data.findIndex((candidate) => candidate.id === item.id))}>
          {item.name}
        </button>
      ),
    },
    { title: '类型', colKey: 'type', cell: ({ row: item }: { row: MockRow }) => <Tag>{item.type}</Tag> },
    { title: '状态', colKey: 'status', cell: ({ row: item }: { row: MockRow }) => <Tag theme={item.status === '正常' || item.status === '运行中' ? 'success' : 'warning'}>{item.status}</Tag> },
    { title: '更新时间', colKey: 'updatedAt' },
  ]

  function act(label: string) {
    setMessage(`${label}: ${row.name} 已进入 Demo 任务队列。生产环境将接入对应 M1/M2/M3 API。`)
  }

  return (
    <div className="workspace-page">
      <div className="workspace-header">
        <div>
          <h1 className="page-title">{viewTitles[view]}</h1>
          <p className="page-subtitle">Demo mock 数据已按生产信息架构组织，实例生命周期仍由“实例工作台”真实 API 执行。</p>
        </div>
        <Space>
          <Button variant="outline" onClick={() => act('刷新')}>刷新</Button>
          <Button onClick={() => act('新建')}>新建</Button>
        </Space>
      </div>

      <div className="mock-metric-grid">
        {mockMetrics(view).map((item) => (
          <Card key={item.label}>
            <div className="mock-metric">
              <span>{item.label}</span>
              <strong>{item.value}</strong>
            </div>
          </Card>
        ))}
      </div>

      <div className="mock-console-grid">
        <Card title="资源列表">
          <Table data={data} columns={columns} rowKey="id" />
        </Card>
        <Card title="详情与操作">
          <div className="mock-detail">
            <h2>{row.name}</h2>
            <dl className="detail-list">
              <dt>ID</dt><dd>{row.id}</dd>
              <dt>类型</dt><dd>{row.type}</dd>
              <dt>状态</dt><dd>{row.status}</dd>
              <dt>摘要</dt><dd>{row.summary}</dd>
            </dl>
            <div className="action-grid">
              {['查看详情', '编辑配置', '发布变更', '查看审计'].map((label) => (
                <Button key={label} variant="outline" onClick={() => act(label)}>{label}</Button>
              ))}
            </div>
            <pre className="mock-output">{message || '选择资源或点击操作后，这里展示 Demo 任务反馈。'}</pre>
          </div>
        </Card>
      </div>
    </div>
  )
}

type MockRow = {
  id: string
  name: string
  type: string
  status: string
  updatedAt: string
  summary: string
}

function mockMetrics(view: ConsoleView) {
  const values: Record<ConsoleView, Array<{ label: string; value: string }>> = {
    'demo-instances': [],
    dashboard: [
      { label: '运行实例', value: '18' },
      { label: 'GPU 使用率', value: '63%' },
      { label: '今日任务', value: '42' },
      { label: '告警', value: '2' },
    ],
    models: [
      { label: '模型', value: '12' },
      { label: '版本', value: '31' },
      { label: '待审核', value: '3' },
      { label: '已发布', value: '9' },
    ],
    kb: [
      { label: '知识库', value: '8' },
      { label: '文档', value: '1.2k' },
      { label: '索引任务', value: '6' },
      { label: '失败', value: '1' },
    ],
    usage: [
      { label: 'Token', value: '18.4M' },
      { label: '推理请求', value: '92k' },
      { label: '存储', value: '3.7TB' },
      { label: '成本', value: '¥8.6k' },
    ],
    settings: [
      { label: '租户', value: '5' },
      { label: '角色', value: '14' },
      { label: 'API Key', value: '27' },
      { label: '策略', value: '19' },
    ],
  }
  return values[view]
}

function mockRows(view: ConsoleView): MockRow[] {
  const rows: Record<ConsoleView, MockRow[]> = {
    'demo-instances': [],
    dashboard: [
      { id: 'dash-prod', name: '生产集群概览', type: 'Dashboard', status: '正常', updatedAt: '2 分钟前', summary: '实例、模型、知识库、任务和告警的统一运营视图。' },
      { id: 'dash-gpu', name: 'GPU 资源看板', type: 'GPU', status: '正常', updatedAt: '5 分钟前', summary: '异构 GPU 库存、分配、温度和利用率趋势。' },
    ],
    models: [
      { id: 'model-qwen', name: 'qwen2.5-72b-instruct', type: 'LLM', status: '运行中', updatedAt: '10 分钟前', summary: '多版本模型资产，后续接入对象存储和推理服务发布链路。' },
      { id: 'model-embed', name: 'bge-large-zh-v1.5', type: 'Embedding', status: '正常', updatedAt: '18 分钟前', summary: '知识库检索使用的向量模型。' },
      { id: 'model-rerank', name: 'bge-reranker-v2', type: 'Rerank', status: '待发布', updatedAt: '1 小时前', summary: 'RAG 排序模型，展示审批和发布动作。' },
    ],
    kb: [
      { id: 'kb-ops', name: '运维知识库', type: 'RAG', status: '正常', updatedAt: '4 分钟前', summary: '包含告警手册、应急预案和运行记录。' },
      { id: 'kb-product', name: '产品文档库', type: 'RAG', status: '索引中', updatedAt: '12 分钟前', summary: '产品设计、API 文档和发布说明。' },
    ],
    usage: [
      { id: 'usage-tenant-a', name: 'tenant-a', type: 'Tenant', status: '正常', updatedAt: '今天', summary: '租户级 token、GPU、存储和推理请求聚合。' },
      { id: 'usage-tenant-b', name: 'tenant-b', type: 'Tenant', status: '正常', updatedAt: '今天', summary: '部门级用量核算和配额趋势。' },
    ],
    settings: [
      { id: 'set-rbac', name: 'RBAC 策略', type: 'Security', status: '正常', updatedAt: '20 分钟前', summary: '租户角色、权限边界和 API Key 管理。' },
      { id: 'set-adapter', name: '组件适配策略', type: 'Architecture', status: '正常', updatedAt: '30 分钟前', summary: '开源组件松耦合、Provider 能力声明和直连例外。' },
    ],
  }
  return rows[view]
}
