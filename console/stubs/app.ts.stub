import { createApp, h, type DefineComponent } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import Layout from './Layout.vue'

const pages = import.meta.glob<DefineComponent>('./Pages/**/*.vue')

createInertiaApp({
  progress: {
    // The delay after which the progress bar will appear, in milliseconds...
    delay: 250,
    // The color of the progress bar...
    color: '#29d',
    // Whether to include the default NProgress styles...
    includeCSS: true,
    // Whether the NProgress spinner will be shown...
    showSpinner: false,
  },
  title: (title) => (title ? `${title} · Goravel Inertia` : 'Goravel Inertia'),
  resolve: (name) => {
    const page = pages[`./Pages/${name}.vue`]
    if (!page) {
      throw new Error(`Page not found: ${name}`)
    }
    return page().then((module) => {
      const component = module.default
      if (!component.layout) {
        component.layout = Layout
      }
      return component
    })
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
