import { createInertiaApp } from '@inertiajs/vue3'
import createServer from '@inertiajs/vue3/server'
import { renderToString } from '@vue/server-renderer'
import { createSSRApp, h, type DefineComponent } from 'vue'
import Layout from './Layout.vue'

const pages = import.meta.glob<DefineComponent>('./Pages/**/*.vue')

// SSR entry. `createServer` starts a Node HTTP server (default :13714, endpoint
// /render) that the Go backend POSTs each Inertia page to.
createServer((page) =>
  createInertiaApp({
    page,
    render: renderToString,
    title: (title) => (title ? `${title} · Goravel Inertia` : 'Goravel Inertia'),
    resolve: (name) => {
      const resolved = pages[`./Pages/${name}.vue`]
      if (!resolved) {
        throw new Error(`Page not found: ${name}`)
      }
      return resolved().then((module) => {
        const component = module.default
        if (!component.layout) {
          component.layout = Layout
        }
        return component
      })
    },
    setup({ App, props, plugin }) {
      return createSSRApp({ render: () => h(App, props) }).use(plugin)
    },
  }),
)
