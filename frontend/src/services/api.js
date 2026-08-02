// API service for communicating with the Go backend
const API_BASE = 'http://localhost:8080/api'

// Helper to get the stored JWT token
function getToken() {
  return localStorage.getItem('token')
}

// Generic fetch wrapper with auth header support
async function request(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  const token = getToken()
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    ...options,
    headers,
  })

  let data
  const responseText = await response.text()
  try {
    data = responseText ? JSON.parse(responseText) : {}
  } catch (err) {
    data = {
      success: false,
      message: responseText || `Server error (${response.status})`,
    }
  }

  return { status: response.status, data }
}

// ========== Auth APIs ==========

export async function registerUser(payload) {
  return request('/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function loginUser(payload) {
  return request('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function getProfile() {
  return request('/auth/profile', {
    method: 'GET',
  })
}

export async function updateUserProfile(payload) {
  return request('/auth/profile', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export async function getUserById(id) {
  return request(`/auth/users/${id}`, {
    method: 'GET',
  })
}

// ========== Provider APIs ==========

export async function createProviderProfile(payload) {
  return request('/provider/profile', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function getProviderProfile() {
  return request('/provider/profile', {
    method: 'GET',
  })
}

export async function updateProviderProfile(payload) {
  return request('/provider/profile', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export async function getAllProviders(page = 1, limit = 12) {
  return request(`/providers?page=${page}&limit=${limit}`, {
    method: 'GET',
  })
}

export async function getProviderById(id) {
  return request(`/providers/${id}`, {
    method: 'GET',
  })
}

// ========== Appointment APIs ==========

export async function createAppointment(payload) {
  return request('/appointments', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function getMyAppointments() {
  return request('/appointments', {
    method: 'GET',
  })
}

export async function updateAppointmentStatus(id, status) {
  return request(`/appointments/${id}/status`, {
    method: 'PATCH',
    body: JSON.stringify({ status }),
  })
}

// ========== Review APIs ==========

export async function createReview(payload) {
  return request('/reviews', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export async function getProviderReviews(providerId) {
  return request(`/reviews/provider/${providerId}`, {
    method: 'GET',
  })
}
