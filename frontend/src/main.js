import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config'

// PrimeVue v3 100% free MIT open-source themes & CSS
import 'primevue/resources/themes/lara-light-green/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
import './style.css'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Dialog from 'primevue/dialog'
import Card from 'primevue/card'
import Rating from 'primevue/rating'
import Tag from 'primevue/tag'
import Calendar from 'primevue/calendar'
import ProgressSpinner from 'primevue/progressspinner'
import Message from 'primevue/message'

const app = createApp(App)

app.use(router)
app.use(PrimeVue, { ripple: true })

app.component('PButton', Button)
app.component('PInputText', InputText)
app.component('PTextarea', Textarea)
app.component('PSelect', Dropdown)
app.component('PDialog', Dialog)
app.component('PCard', Card)
app.component('PRating', Rating)
app.component('PTag', Tag)
app.component('PCalendar', Calendar)
app.component('PProgressSpinner', ProgressSpinner)
app.component('PMessage', Message)

app.mount('#app')
