# 🎨 Guía de Integración Frontend (Vue.js)

## 📦 Componente Vue para Términos y Condiciones

### 1. Crear Vista de Términos

Crear archivo: `src/views/TermsAcceptance.vue`

```vue
<template>
  <div class="terms-container">
    <!-- Loading State -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Cargando términos y condiciones...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h1>Error</h1>
      <p>{{ error }}</p>
      <button @click="reloadPage" class="btn-secondary">Reintentar</button>
    </div>

    <!-- Expired State -->
    <div v-else-if="status === 'EXPIRED'" class="expired-state">
      <div class="icon">⏰</div>
      <h1>Link Expirado</h1>
      <p>Este link de términos y condiciones ha expirado.</p>
      <p class="help-text">Por favor, solicita un nuevo link a través del chat.</p>
    </div>

    <!-- Accepted State -->
    <div v-else-if="status === 'ACCEPTED'" class="accepted-state">
      <div class="icon success">✓</div>
      <h1>Términos Aceptados</h1>
      <p>Ya has aceptado los términos y condiciones.</p>
      <p class="timestamp">Aceptado el: {{ formatDate(acceptedAt) }}</p>
      <p class="help-text">Puedes cerrar esta ventana.</p>
    </div>

    <!-- Rejected State -->
    <div v-else-if="status === 'REJECTED'" class="rejected-state">
      <div class="icon">✗</div>
      <h1>Términos Rechazados</h1>
      <p>Has rechazado los términos y condiciones.</p>
      <p class="timestamp">Rechazado el: {{ formatDate(rejectedAt) }}</p>
    </div>

    <!-- Pending State - Show Terms -->
    <div v-else-if="status === 'PENDING'" class="pending-state">
      <div class="terms-header">
        <h1>Términos y Condiciones de Entrega</h1>
        <p class="expiry-info">
          Este link expira el: <strong>{{ formatDate(expiresAt) }}</strong>
        </p>
      </div>

      <div class="terms-content">
        <div class="terms-scroll">
          <h2>1. Aceptación de Términos</h2>
          <p>
            Al aceptar estos términos, usted confirma que ha recibido la entrega
            de los dispensers en las condiciones acordadas.
          </p>

          <h2>2. Condiciones de Recepción</h2>
          <ul>
            <li>Los dispensers han sido entregados en perfecto estado</li>
            <li>La cantidad recibida coincide con lo solicitado</li>
            <li>Los números de serie han sido verificados</li>
          </ul>

          <h2>3. Responsabilidad</h2>
          <p>
            A partir de la aceptación, usted se hace responsable del cuidado
            y mantenimiento de los equipos entregados.
          </p>

          <h2>4. Garantía</h2>
          <p>
            Los equipos cuentan con garantía según los términos establecidos
            en el contrato de servicio.
          </p>

          <!-- Agregar más términos según necesites -->
        </div>

        <div class="terms-checkbox">
          <label>
            <input type="checkbox" v-model="agreedToTerms" />
            He leído y acepto los términos y condiciones
          </label>
        </div>
      </div>

      <div class="terms-actions">
        <button 
          @click="acceptTerms" 
          :disabled="!agreedToTerms || processing"
          class="btn-primary"
        >
          <span v-if="!processing">Aceptar Términos</span>
          <span v-else>Procesando...</span>
        </button>
        <button 
          @click="showRejectConfirm = true"
          :disabled="processing"
          class="btn-secondary"
        >
          Rechazar
        </button>
      </div>
    </div>

    <!-- Reject Confirmation Modal -->
    <div v-if="showRejectConfirm" class="modal-overlay" @click="showRejectConfirm = false">
      <div class="modal-content" @click.stop>
        <h2>¿Estás seguro?</h2>
        <p>¿Realmente deseas rechazar los términos y condiciones?</p>
        <div class="modal-actions">
          <button @click="rejectTerms" class="btn-danger">Sí, Rechazar</button>
          <button @click="showRejectConfirm = false" class="btn-secondary">Cancelar</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const token = route.params.token

// State
const status = ref('PENDING')
const loading = ref(true)
const error = ref(null)
const processing = ref(false)
const agreedToTerms = ref(false)
const showRejectConfirm = ref(false)
const expiresAt = ref(null)
const acceptedAt = ref(null)
const rejectedAt = ref(null)

// API Configuration
const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// Load terms status on mount
onMounted(async () => {
  await loadTermsStatus()
})

// Load terms status from API
const loadTermsStatus = async () => {
  loading.value = true
  error.value = null

  try {
    const { data } = await axios.get(`${API_BASE}/terms/${token}`)
    status.value = data.status
    expiresAt.value = data.expiresAt
    acceptedAt.value = data.acceptedAt
    rejectedAt.value = data.rejectedAt
  } catch (err) {
    console.error('Error loading terms:', err)
    if (err.response?.status === 404) {
      error.value = 'El link proporcionado no es válido o ha expirado.'
    } else {
      error.value = 'Error al cargar los términos. Por favor, intenta nuevamente.'
    }
  } finally {
    loading.value = false
  }
}

// Accept terms
const acceptTerms = async () => {
  if (!agreedToTerms.value) return

  processing.value = true
  error.value = null

  try {
    const { data } = await axios.post(`${API_BASE}/terms/${token}/accept`)
    
    status.value = data.status
    acceptedAt.value = data.acceptedAt
    
    // Mostrar mensaje de éxito
    alert('✓ Términos aceptados exitosamente')
    
  } catch (err) {
    console.error('Error accepting terms:', err)
    
    if (err.response?.status === 410) {
      error.value = 'El link ha expirado. Por favor, solicita uno nuevo.'
      status.value = 'EXPIRED'
    } else {
      error.value = err.response?.data?.error || 'Error al aceptar términos. Intenta nuevamente.'
    }
  } finally {
    processing.value = false
  }
}

// Reject terms
const rejectTerms = async () => {
  processing.value = true
  showRejectConfirm.value = false
  error.value = null

  try {
    const { data } = await axios.post(`${API_BASE}/terms/${token}/reject`)
    
    status.value = data.status
    rejectedAt.value = data.rejectedAt
    
    alert('Términos rechazados')
    
  } catch (err) {
    console.error('Error rejecting terms:', err)
    error.value = err.response?.data?.error || 'Error al rechazar términos. Intenta nuevamente.'
  } finally {
    processing.value = false
  }
}

// Format date for display
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Reload page
const reloadPage = () => {
  window.location.reload()
}
</script>

<style scoped>
.terms-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

/* Loading State */
.loading {
  text-align: center;
  padding: 4rem 2rem;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Status States */
.error-state,
.expired-state,
.accepted-state,
.rejected-state {
  text-align: center;
  padding: 3rem 2rem;
}

.icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

.icon.success {
  color: #27ae60;
}

.error-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

h1 {
  color: #2c3e50;
  margin-bottom: 1rem;
}

.timestamp {
  color: #7f8c8d;
  font-size: 0.9rem;
  margin-top: 1rem;
}

.help-text {
  color: #95a5a6;
  font-size: 0.9rem;
  margin-top: 1rem;
}

/* Pending State - Terms Display */
.pending-state {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.terms-header {
  text-align: center;
  border-bottom: 2px solid #ecf0f1;
  padding-bottom: 1.5rem;
}

.expiry-info {
  color: #7f8c8d;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.terms-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.terms-scroll {
  max-height: 400px;
  overflow-y: auto;
  padding: 1.5rem;
  border: 1px solid #ecf0f1;
  border-radius: 8px;
  background-color: #f8f9fa;
}

.terms-scroll h2 {
  color: #2c3e50;
  font-size: 1.2rem;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
}

.terms-scroll h2:first-child {
  margin-top: 0;
}

.terms-scroll ul {
  margin-left: 1.5rem;
}

.terms-scroll li {
  margin-bottom: 0.5rem;
  line-height: 1.6;
}

.terms-checkbox {
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 2px solid #dee2e6;
}

.terms-checkbox label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  font-weight: 500;
}

.terms-checkbox input[type="checkbox"] {
  width: 20px;
  height: 20px;
  cursor: pointer;
}

.terms-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

/* Buttons */
.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background-color: #27ae60;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #229954;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(39, 174, 96, 0.3);
}

.btn-primary:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
  opacity: 0.6;
}

.btn-secondary {
  background-color: #ecf0f1;
  color: #2c3e50;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #d5dbdb;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.btn-danger:hover {
  background-color: #c0392b;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  max-width: 400px;
  width: 90%;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.modal-content h2 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.modal-content p {
  margin-bottom: 1.5rem;
  color: #7f8c8d;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

/* Responsive */
@media (max-width: 768px) {
  .terms-container {
    padding: 1rem;
  }

  .terms-actions {
    flex-direction: column;
  }

  .btn-primary,
  .btn-secondary,
  .btn-danger {
    width: 100%;
  }

  .modal-actions {
    flex-direction: column;
  }
}
</style>
```

---

## 📍 2. Configurar Router

En `src/router/index.js`:

```javascript
import { createRouter, createWebHistory } from 'vue-router'
import TermsAcceptance from '@/views/TermsAcceptance.vue'

const routes = [
  // ... otras rutas
  {
    path: '/terms/:token',
    name: 'TermsAcceptance',
    component: TermsAcceptance,
    meta: {
      title: 'Términos y Condiciones'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

---

## 🔧 3. Configurar Variables de Entorno

Crear `.env` o `.env.local`:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Para producción (`.env.production`):

```env
VITE_API_BASE_URL=https://api.mi-dominio.com/api/v1
```

---

## 📱 4. Componente Alternativo Simplificado

Si prefieres algo más simple:

```vue
<template>
  <div class="terms-simple">
    <div v-if="loading">Cargando...</div>
    
    <div v-else-if="status === 'PENDING'" class="terms-pending">
      <h1>Términos y Condiciones</h1>
      <div class="terms-text">
        <!-- Contenido de términos -->
      </div>
      <button @click="accept">Aceptar</button>
      <button @click="reject">Rechazar</button>
    </div>
    
    <div v-else-if="status === 'ACCEPTED'">
      ✓ Términos aceptados
    </div>
    
    <div v-else>
      Estado: {{ status }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const token = route.params.token
const API_BASE = 'http://localhost:8080/api/v1'

const loading = ref(true)
const status = ref('PENDING')

onMounted(async () => {
  const { data } = await axios.get(`${API_BASE}/terms/${token}`)
  status.value = data.status
  loading.value = false
})

const accept = async () => {
  await axios.post(`${API_BASE}/terms/${token}/accept`)
  status.value = 'ACCEPTED'
}

const reject = async () => {
  await axios.post(`${API_BASE}/terms/${token}/reject`)
  status.value = 'REJECTED'
}
</script>
```

---

## 🧪 5. Testing del Componente

Crear archivo: `src/views/__tests__/TermsAcceptance.spec.js`

```javascript
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import axios from 'axios'
import TermsAcceptance from '../TermsAcceptance.vue'

vi.mock('axios')

describe('TermsAcceptance', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('loads and displays pending terms', async () => {
    axios.get.mockResolvedValueOnce({
      data: {
        status: 'PENDING',
        expiresAt: '2025-12-26T10:00:00Z'
      }
    })

    const wrapper = mount(TermsAcceptance, {
      global: {
        mocks: {
          $route: { params: { token: 'test-token' } }
        }
      }
    })

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Términos y Condiciones')
  })

  it('accepts terms successfully', async () => {
    axios.get.mockResolvedValueOnce({
      data: { status: 'PENDING', expiresAt: '2025-12-26T10:00:00Z' }
    })
    
    axios.post.mockResolvedValueOnce({
      data: { status: 'ACCEPTED', message: 'Success' }
    })

    const wrapper = mount(TermsAcceptance, {
      global: {
        mocks: {
          $route: { params: { token: 'test-token' } }
        }
      }
    })

    await wrapper.vm.$nextTick()
    
    // Simulate accept
    await wrapper.vm.acceptTerms()

    expect(axios.post).toHaveBeenCalledWith(
      expect.stringContaining('/accept')
    )
  })
})
```

---

## 🎨 6. Composable Reutilizable

Crear `src/composables/useTermsSession.js`:

```javascript
import { ref } from 'vue'
import axios from 'axios'

export function useTermsSession(token) {
  const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  
  const status = ref('PENDING')
  const loading = ref(false)
  const error = ref(null)
  const expiresAt = ref(null)
  const acceptedAt = ref(null)
  const rejectedAt = ref(null)

  const loadStatus = async () => {
    loading.value = true
    error.value = null
    
    try {
      const { data } = await axios.get(`${API_BASE}/terms/${token}`)
      status.value = data.status
      expiresAt.value = data.expiresAt
      acceptedAt.value = data.acceptedAt
      rejectedAt.value = data.rejectedAt
    } catch (err) {
      error.value = err.response?.data?.error || 'Error loading terms'
      throw err
    } finally {
      loading.value = false
    }
  }

  const accept = async () => {
    loading.value = true
    error.value = null
    
    try {
      const { data } = await axios.post(`${API_BASE}/terms/${token}/accept`)
      status.value = data.status
      acceptedAt.value = data.acceptedAt
      return data
    } catch (err) {
      error.value = err.response?.data?.error || 'Error accepting terms'
      throw err
    } finally {
      loading.value = false
    }
  }

  const reject = async () => {
    loading.value = true
    error.value = null
    
    try {
      const { data } = await axios.post(`${API_BASE}/terms/${token}/reject`)
      status.value = data.status
      rejectedAt.value = data.rejectedAt
      return data
    } catch (err) {
      error.value = err.response?.data?.error || 'Error rejecting terms'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    status,
    loading,
    error,
    expiresAt,
    acceptedAt,
    rejectedAt,
    loadStatus,
    accept,
    reject
  }
}
```

Uso del composable:

```vue
<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useTermsSession } from '@/composables/useTermsSession'

const route = useRoute()
const token = route.params.token

const {
  status,
  loading,
  error,
  expiresAt,
  acceptedAt,
  loadStatus,
  accept,
  reject
} = useTermsSession(token)

onMounted(() => {
  loadStatus()
})
</script>
```

---

## 🔔 7. Notificaciones con Vue Toastification

Instalar:
```bash
npm install vue-toastification
```

Configurar en `main.js`:
```javascript
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

app.use(Toast)
```

Usar en componente:
```javascript
import { useToast } from 'vue-toastification'

const toast = useToast()

const acceptTerms = async () => {
  try {
    await accept()
    toast.success('Términos aceptados exitosamente')
  } catch (err) {
    toast.error('Error al aceptar términos')
  }
}
```

---

## 📱 8. Ejemplo Responsive Mejorado

```vue
<style scoped>
@media (max-width: 768px) {
  .terms-container {
    padding: 1rem;
  }
  
  .terms-header h1 {
    font-size: 1.5rem;
  }
  
  .terms-scroll {
    max-height: 300px;
  }
  
  .terms-actions {
    flex-direction: column;
  }
  
  button {
    width: 100%;
  }
}

@media (prefers-color-scheme: dark) {
  .terms-container {
    background-color: #1a1a1a;
    color: #e0e0e0;
  }
  
  .terms-scroll {
    background-color: #2d2d2d;
    border-color: #444;
  }
}
</style>
```

---

## ✅ Checklist de Integración Frontend

- [ ] Componente `TermsAcceptance.vue` creado
- [ ] Ruta `/terms/:token` configurada en router
- [ ] Variable `VITE_API_BASE_URL` configurada
- [ ] Axios instalado y configurado
- [ ] CORS configurado en backend
- [ ] Pruebas con token real realizadas
- [ ] Estilos responsive verificados
- [ ] Manejo de errores implementado
- [ ] Loading states añadidos
- [ ] Notificaciones de éxito/error configuradas

---

**¡Listo para integrar! 🎉**
