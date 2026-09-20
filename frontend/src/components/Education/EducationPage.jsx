import { useState } from 'react';
import { faqItems, ispuGuide, emergencyContacts } from '../../data/education';

export default function EducationPage({ onClose }) {
  const [openFaq, setOpenFaq] = useState(null);

  return (
    <main className="app-page education-page">
      <div className="page-shell">
        <button className="page-back" onClick={onClose}>Kembali ke dashboard</button>
        <header className="page-hero education-hero">
          <div>
            <span className="page-eyebrow">Pusat pengetahuan</span>
            <h1>Edukasi karhutla & kualitas udara</h1>
            <p>Pahami arti status udara, kenali risikonya, dan pilih tindakan yang tepat untuk melindungi diri serta keluarga.</p>
          </div>
          <div className="hero-stat"><strong>5</strong><span>level ISPU</span></div>
        </header>

        <section className="page-section education-section">
          <div className="section-heading"><span className="page-eyebrow">Panduan kesehatan</span><h2>Setiap level punya tindakan berbeda</h2></div>
          <div className="ispu-guide">
            {ispuGuide.map((item) => (
              <div className="ispu-card" key={item.category} style={{ borderLeftColor: item.color }}>
                <div className="ispu-card-header">
                  <span className="ispu-range">{item.range}</span>
                  <strong style={{ color: item.color }}>{item.category}</strong>
                </div>
                <p><strong>Aktivitas:</strong> {item.activity}</p>
                <p><strong>Masker:</strong> {item.mask}</p>
                <p><strong>Tindakan medis:</strong> {item.doctor}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="education-lower-grid">
        <section className="page-section education-section education-faq-section">
          <div className="section-heading"><span className="page-eyebrow">Jawaban singkat</span><h2>Pertanyaan umum</h2></div>
          <div className="faq-list faq-list--card">
            {faqItems.map((item, idx) => (
              <div className="faq-item" key={idx}>
                <button
                  className="faq-question"
                  onClick={() => setOpenFaq(openFaq === idx ? null : idx)}
                >
                  {item.q}
                  <span>{openFaq === idx ? '−' : '+'}</span>
                </button>
                {openFaq === idx && <p className="faq-answer">{item.a}</p>}
              </div>
            ))}
          </div>
        </section>

        <section className="page-section education-section education-contact-section">
          <div className="section-heading"><span className="page-eyebrow">Butuh bantuan segera?</span><h2>Layanan darurat</h2></div>
          <div className="contact-list">
            {emergencyContacts.map((c) => (
              <a className="contact-item" href={`tel:${c.number}`} key={c.number}>
                <span>{c.name}</span>
                <strong>{c.number}</strong>
              </a>
            ))}
          </div>
        </section>
        </div>
      </div>
    </main>
  );
}