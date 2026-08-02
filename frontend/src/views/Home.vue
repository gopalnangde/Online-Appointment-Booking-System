<template>
  <div class="landing-wrapper">
    <div class="container">
      <!-- Hero Banner Card -->
      <section class="hero-banner-card">
        <div class="hero-content">
          <div class="badge-pill">
            <i class="pi pi-sparkles"></i>
            <span>#1 Online Appointment Platform</span>
          </div>

          <h1 class="hero-title">
            Book Appointments <span class="highlight">Effortlessly</span>
          </h1>

          <p class="hero-subtitle">
            Connect with verified service providers, explore expert specializations,
            and schedule your appointments in seconds.
          </p>

          <div class="hero-actions">
            <router-link to="/register">
              <PButton label="Get Started Free" icon="pi pi-user-plus" class="p-button-lg hero-btn-white" />
            </router-link>
            <router-link to="/providers">
              <PButton label="Browse Providers" icon="pi pi-search" class="p-button-lg hero-btn-glass" />
            </router-link>
          </div>
        </div>
      </section>

      <!-- Stats Bar Row -->
      <section class="stats-row">
        <div class="stat-card">
          <div class="stat-icon-bg">
            <i class="pi pi-users"></i>
          </div>
          <div>
            <h3 class="stat-number">500+</h3>
            <p class="stat-label">Service Providers</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon-bg">
            <i class="pi pi-calendar-check"></i>
          </div>
          <div>
            <h3 class="stat-number">10K+</h3>
            <p class="stat-label">Appointments Booked</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon-bg">
            <i class="pi pi-briefcase"></i>
          </div>
          <div>
            <h3 class="stat-number">50+</h3>
            <p class="stat-label">Categories & Specialties</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon-bg star-icon-bg">
            <i class="pi pi-star-fill"></i>
          </div>
          <div>
            <h3 class="stat-number">4.9 ★</h3>
            <p class="stat-label">User Satisfaction</p>
          </div>
        </div>
      </section>

      <!-- Features Section -->
      <section class="features-section">
        <div class="section-header text-center">
          <span class="sub-tag">Designed for Simplicity</span>
          <h2 class="section-title">Why Choose Appointly?</h2>
          <p class="section-subtitle">
            Everything you need for seamless appointment management and client scheduling.
          </p>
        </div>

        <div class="grid grid-3">
          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-search"></i>
            </div>
            <h3>Discover Experts</h3>
            <p>Browse verified healthcare and professional service providers with clear experience details and locations.</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-calendar"></i>
            </div>
            <h3>Instant Booking</h3>
            <p>Select your date and time slot in real-time — no phone calls or manual waiting required.</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-sliders-h"></i>
            </div>
            <h3>Smart Status Tracking</h3>
            <p>Track your appointment status through clear phases: Pending ➔ Confirmed ➔ Completed.</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-star-fill"></i>
            </div>
            <h3>Ratings & Feedback</h3>
            <p>Leave reviews and star ratings after completed appointments to maintain high service quality.</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-id-card"></i>
            </div>
            <h3>Provider Practice Hub</h3>
            <p>Dedicated dashboard for service providers to build their practice profile and manage bookings.</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon">
              <i class="pi pi-shield"></i>
            </div>
            <h3>Secure & Private</h3>
            <p>Built with enterprise-grade JWT authentication and role-based access control.</p>
          </div>
        </div>
      </section>

      <!-- How it Works Steps -->
      <section class="steps-section">
        <div class="section-header text-center">
          <span class="sub-tag">Three Simple Steps</span>
          <h2 class="section-title">How Appointly Works</h2>
        </div>

        <div class="steps-grid">
          <div class="step-card">
            <div class="step-badge">1</div>
            <h3>Explore & Choose</h3>
            <p>Search for service providers by location, category, or specialization.</p>
          </div>

          <div class="step-card">
            <div class="step-badge">2</div>
            <h3>Pick Date & Time</h3>
            <p>Select your preferred appointment date and time slot with instant confirmation.</p>
          </div>

          <div class="step-card">
            <div class="step-badge">3</div>
            <h3>Attend & Review</h3>
            <p>Connect with your provider and leave ratings after your appointment.</p>
          </div>
        </div>
      </section>

      <!-- Bottom CTA Card -->
      <section class="cta-banner-card">
        <h2>Ready to Experience Seamless Bookings?</h2>
        <p>Join thousands of customers and providers simplifying appointment scheduling today.</p>
        <div class="cta-actions">
          <router-link to="/register">
            <PButton label="Create Account Now" icon="pi pi-arrow-right" iconPos="right" class="p-button-lg p-button-primary" />
          </router-link>
          <router-link to="/login">
            <PButton label="Sign In" icon="pi pi-sign-in" class="p-button-lg p-button-outlined p-button-secondary" />
          </router-link>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
}
</script>

<style scoped>
.landing-wrapper {
  padding: 10px 0 60px;
}

/* Hero Banner Card */
.hero-banner-card {
  background: linear-gradient(135deg, #0052CC 0%, #0066FF 50%, #2563EB 100%);
  border-radius: var(--radius-xl);
  padding: 60px 48px;
  color: var(--color-white);
  box-shadow: 0 12px 36px rgba(0, 102, 255, 0.2);
  margin-bottom: 32px;
  position: relative;
  overflow: hidden;
}

.hero-content {
  max-width: 720px;
}

.badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(8px);
  padding: 6px 16px;
  border-radius: var(--radius-full);
  font-size: 0.85rem;
  font-weight: 700;
  margin-bottom: 24px;
}

.hero-title {
  font-size: 3rem;
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.03em;
  margin-bottom: 20px;
}

.hero-title .highlight {
  color: #60A5FA;
}

.hero-subtitle {
  font-size: 1.15rem;
  opacity: 0.92;
  line-height: 1.6;
  margin-bottom: 36px;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

/* White & Glass buttons inside hero */
.hero-btn-white {
  background: #FFFFFF !important;
  color: #0066FF !important;
  border-color: #FFFFFF !important;
}

.hero-btn-white:hover {
  background: #F0F6FF !important;
  color: #0052CC !important;
}

.hero-btn-glass {
  background: rgba(255, 255, 255, 0.15) !important;
  color: #FFFFFF !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  backdrop-filter: blur(6px);
}

.hero-btn-glass:hover {
  background: rgba(255, 255, 255, 0.25) !important;
}

/* Stats Row */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 60px;
}

@media (max-width: 900px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 500px) {
  .stats-row {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 24px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: var(--transition);
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.stat-icon-bg {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: var(--color-primary-pale);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  flex-shrink: 0;
}

.star-icon-bg {
  background: #FEF3C7;
  color: #F59E0B;
}

.stat-number {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--color-dark);
}

.stat-label {
  font-size: 0.85rem;
  color: var(--color-gray-500);
  font-weight: 600;
}

/* Features Section */
.features-section {
  margin-bottom: 60px;
}

.text-center {
  text-align: center;
}

.section-header {
  max-width: 600px;
  margin: 0 auto 40px;
}

.sub-tag {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--color-primary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.section-title {
  font-size: 2rem;
  font-weight: 800;
  color: var(--color-dark);
  margin: 6px 0 10px;
}

.section-subtitle {
  font-size: 0.95rem;
  color: var(--color-gray-500);
}

.feature-card {
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 28px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-sm);
  transition: var(--transition);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary-light);
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: var(--color-primary-pale);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
}

.feature-card h3 {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--color-dark);
}

.feature-card p {
  font-size: 0.9rem;
  color: var(--color-gray-500);
  line-height: 1.6;
}

/* Steps Section */
.steps-section {
  margin-bottom: 60px;
}

.steps-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

@media (max-width: 768px) {
  .steps-grid {
    grid-template-columns: 1fr;
  }
}

.step-card {
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 32px 24px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-sm);
  text-align: center;
  position: relative;
}

.step-badge {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: var(--color-white);
  font-weight: 800;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  box-shadow: 0 4px 12px rgba(0, 102, 255, 0.25);
}

.step-card h3 {
  font-size: 1.1rem;
  font-weight: 700;
  margin-bottom: 8px;
}

.step-card p {
  font-size: 0.9rem;
  color: var(--color-gray-500);
}

/* CTA Banner Card */
.cta-banner-card {
  background: var(--color-white);
  border-radius: var(--radius-xl);
  padding: 48px 36px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-md);
  text-align: center;
}

.cta-banner-card h2 {
  font-size: 2rem;
  font-weight: 800;
  color: var(--color-dark);
}

.cta-banner-card p {
  color: var(--color-gray-500);
  margin: 10px 0 28px;
  font-size: 1rem;
}

.cta-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}
</style>
