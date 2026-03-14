type WailsRuntime = {
  Call: {
    ByName: <T = unknown>(methodName: string, ...args: unknown[]) => Promise<T>
  }
}

type LegacyBindingTree = Record<string, Record<string, (...args: unknown[]) => Promise<unknown>>>

const bindingMethods = {
  SkillBinding: {
    ListInstalled: 'skillmanager/internal/binding.SkillBinding.ListInstalled',
    GetDetail: 'skillmanager/internal/binding.SkillBinding.GetDetail',
    Install: 'skillmanager/internal/binding.SkillBinding.Install',
    Uninstall: 'skillmanager/internal/binding.SkillBinding.Uninstall',
    Update: 'skillmanager/internal/binding.SkillBinding.Update',
    UpdateContent: 'skillmanager/internal/binding.SkillBinding.UpdateContent',
    AssignAgents: 'skillmanager/internal/binding.SkillBinding.AssignAgents'
  },
  AgentBinding: {
    ListAgents: 'skillmanager/internal/binding.AgentBinding.ListAgents',
    DetectInstalled: 'skillmanager/internal/binding.AgentBinding.DetectInstalled',
    AddCustomAgent: 'skillmanager/internal/binding.AgentBinding.AddCustomAgent',
    RemoveAgent: 'skillmanager/internal/binding.AgentBinding.RemoveAgent',
    ToggleAgent: 'skillmanager/internal/binding.AgentBinding.ToggleAgent'
  },
  RegistryBinding: {
    ListRegistries: 'skillmanager/internal/binding.RegistryBinding.ListRegistries',
    AddRegistry: 'skillmanager/internal/binding.RegistryBinding.AddRegistry',
    RemoveRegistry: 'skillmanager/internal/binding.RegistryBinding.RemoveRegistry',
    Browse: 'skillmanager/internal/binding.RegistryBinding.Browse',
    Search: 'skillmanager/internal/binding.RegistryBinding.Search'
  },
  ConfigBinding: {
    GetConfig: 'skillmanager/internal/binding.ConfigBinding.GetConfig',
    UpdateProxy: 'skillmanager/internal/binding.ConfigBinding.UpdateProxy',
    GetProxy: 'skillmanager/internal/binding.ConfigBinding.GetProxy',
    SetProxy: 'skillmanager/internal/binding.ConfigBinding.SetProxy'
  },
  InventoryBinding: {
    BuildReport: 'skillmanager/internal/binding.InventoryBinding.BuildReport'
  }
} as const

type WindowWithWails = Window & {
  wails?: WailsRuntime
}

let runtimePromise: Promise<WailsRuntime> | null = null

function getWindow(): WindowWithWails {
  return window as WindowWithWails
}

function loadRuntimeScript(): Promise<void> {
  return new Promise((resolve, reject) => {
    const existingScript = document.querySelector<HTMLScriptElement>('script[data-wails-runtime="true"]')
    if (existingScript) {
      if (getWindow().wails?.Call?.ByName) {
        resolve()
        return
      }
      existingScript.addEventListener('load', () => resolve(), { once: true })
      existingScript.addEventListener('error', () => reject(new Error('Failed to load /wails/runtime.js')), { once: true })
      return
    }

    const script = document.createElement('script')
    script.type = 'module'
    script.src = '/wails/runtime.js'
    script.dataset.wailsRuntime = 'true'
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Failed to load /wails/runtime.js'))
    document.head.appendChild(script)
  })
}

function summarizeResult(result: unknown): unknown {
  if (Array.isArray(result)) {
    return { type: 'array', length: result.length }
  }
  if (result && typeof result === 'object') {
    return {
      type: 'object',
      keys: Object.keys(result as Record<string, unknown>).slice(0, 8)
    }
  }
  return result
}

async function getRuntime(): Promise<WailsRuntime> {
  const currentWindow = getWindow()
  if (currentWindow.wails?.Call?.ByName) {
    return currentWindow.wails
  }

  if (!runtimePromise) {
    console.info('[wails] loading runtime bridge')
    runtimePromise = loadRuntimeScript()
      .then(() => {
        const runtime = getWindow().wails
        if (!runtime?.Call?.ByName) {
          throw new Error('Wails runtime loaded, but Call.ByName is unavailable')
        }
        console.info('[wails] runtime bridge ready')
        return runtime
      })
      .catch((error) => {
        runtimePromise = null
        throw error
      })
  }

  return runtimePromise
}

async function callBinding<T>(methodName: string, ...args: unknown[]): Promise<T> {
  const runtime = await getRuntime()
  console.info('[wails] calling', methodName, args)
  try {
    const result = await runtime.Call.ByName<T>(methodName, ...args)
    console.info('[wails] resolved', methodName, summarizeResult(result))
    return result
  } catch (error) {
    console.error('[wails] call failed', methodName, error)
    throw error
  }
}

function buildLegacyBridge(): LegacyBindingTree {
  const app: LegacyBindingTree = {}

  for (const [bindingName, methods] of Object.entries(bindingMethods)) {
    app[bindingName] = {}
    for (const [methodName, fqMethodName] of Object.entries(methods)) {
      app[bindingName][methodName] = (...args: unknown[]) => callBinding(fqMethodName, ...args)
    }
  }

  return app
}

export async function bootstrapWailsCompat(): Promise<void> {
  const currentWindow = getWindow()
  const bridgeWindow = currentWindow as unknown as {
    go?: {
      main?: {
        App?: LegacyBindingTree
      }
    }
  }

  if (bridgeWindow.go?.main?.App) {
    console.info('[wails] legacy window.go bridge already available')
    return
  }

  await getRuntime()

  if (!bridgeWindow.go) {
    bridgeWindow.go = {}
  }
  if (!bridgeWindow.go.main) {
    bridgeWindow.go.main = {}
  }
  bridgeWindow.go.main.App = buildLegacyBridge()
  console.info('[wails] installed window.go.main.App compatibility bridge')
}
