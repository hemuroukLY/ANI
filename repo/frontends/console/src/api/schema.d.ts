export type paths = {
  '/demo/instances': {
    get: {
      responses: {
        200: {
          content: {
            'application/json': {
              items: DemoInstance[]
              total: number
            }
          }
        }
      }
    }
    post: {
      requestBody: {
        content: {
          'application/json': DemoCreateInstanceRequest
        }
      }
      responses: {
        201: {
          content: {
            'application/json': DemoCreateInstanceResponse
          }
        }
      }
    }
  }
  '/demo/instances/{instance_id}/lifecycle': {
    post: {
      parameters: {
        path: { instance_id: string }
      }
      requestBody: {
        content: {
          'application/json': { action: 'start' | 'stop' | 'restart' | 'resize' | 'delete'; cpu?: string; memory?: string }
        }
      }
      responses: {
        200: {
          content: {
            'application/json': DemoInstance
          }
        }
      }
    }
  }
  '/demo/instances/{instance_id}': {
    get: {
      parameters: {
        path: { instance_id: string }
      }
      responses: {
        200: {
          content: {
            'application/json': DemoInstance
          }
        }
      }
    }
  }
  '/demo/instances/{instance_id}/ops/{action}': {
    get: {
      parameters: {
        path: { instance_id: string; action: 'logs' | 'events' | 'metrics' | 'terminal' | 'exec' }
      }
      responses: {
        200: {
          content: {
            'application/json': DemoOpsResult
          }
        }
      }
    }
  }
  '/demo/instances/{instance_id}/console': {
    post: {
      parameters: {
        path: { instance_id: string }
      }
      requestBody: {
        content: {
          'application/json': { protocol: string }
        }
      }
      responses: {
        200: {
          content: {
            'application/json': DemoOpsResult
          }
        }
      }
    }
  }
  '/demo/instances/{instance_id}/console/exec': {
    post: {
      parameters: {
        path: { instance_id: string }
      }
      requestBody: {
        content: {
          'application/json': { command: string }
        }
      }
      responses: {
        200: {
          content: {
            'application/json': { command: string; output: string; exit_code: number; cwd: string }
          }
        }
      }
    }
  }
  '/models': {
    get: {
      responses: {
        200: {
          content: {
            'application/json': { items: Array<Record<string, unknown>> }
          }
        }
      }
    }
  }
  '/inference-services': {
    get: {
      parameters: {
        query?: { status?: string }
      }
      responses: {
        200: {
          content: {
            'application/json': { items: Array<Record<string, unknown>> }
          }
        }
      }
    }
  }
}

export type DemoInstance = {
  id: string
  tenant_id: string
  name: string
  kind: 'vm' | 'container' | 'gpu_container'
  status: string
  provider: string
  resource_refs: string[]
  endpoint: string
  created_at: string
  updated_at: string
}

export type DemoCreateInstanceRequest = {
  kind: 'vm' | 'container' | 'gpu_container'
  name: string
  cpu?: string
  memory?: string
  boot_image?: string
  image?: string
  gpu_vendor?: string
  gpu_model?: string
  gpu_count?: number
  auto_start?: boolean
  description?: string
}

export type DemoCreateInstanceResponse = {
  instance: DemoInstance
  audit_id: string
  manifests: Array<{ name: string; kind: string; provider: string; content: string }>
  timeline: Array<{ name: string; status: string; detail: string }>
  demo_notice: string
}

export type DemoOpsResult = {
  action: string
  accepted: boolean
  session_id: string
  protocol: string
  connect_url: string
  output: string
  reason: string
  expires_at: string
}
