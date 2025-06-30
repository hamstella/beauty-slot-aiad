# Mobile‑First Beauty Salon Booking App – UI/UX Proposal

## Goal
> **Enable customers to complete a reservation with just a few taps, using one hand, on any smartphone browser.**

---

## 1. Information Architecture (IA)

```
┌── Home (Reservation list & “+” FAB) ────────────────┐
│  Recent reservations – 1‑tap rebook                │
│                                                    │
│       ＋  New reservation (FAB)                    │
└────────────────────────────────────────────────────┘
            │ tap
            ▼
┌── Step 1. Service Select (Bottom Sheet) ───────────┐
│  Cut / Color / Perm … (card grid)                  │
└──┬─────────────────────────────────────────────────┘
   │ auto‑save, swipe‑to‑close
   ▼
┌── Step 2. Date & Time (Bottom Sheet) ──────────────┐
│  Horizontal calendar + available time chips        │
└──┬─────────────────────────────────────────────────┘
   │
   ▼
┌── Step 3. Stylist (Optional) ──────────────────────┐
│  “No preference“ / “Same as last time” / list      │
└──┬─────────────────────────────────────────────────┘
   │
   ▼
┌── Step 4. Confirm & Done (Dialog/Sheet) ───────────┐
│  Date‑time · Service · Stylist · Total price       │
│  [Reserve] (primary)                               │
└────────────────────────────────────────────────────┘
```

* **One purpose per screen**; swipe‑back / back‑button moves one level up  
* Each step is a **Bottom Sheet**, layering over the previous screen to preserve spatial context  
* Auto‑save allows interruption & resume at the same step  

---

## 2. Key Screens & UI Details

| Screen | Main UI elements | UX points |
| ------ | ---------------- | --------- |
| **Home** | • 3 most‑recent reservations (cards)<br>• **FAB** “＋ Reserve” | Re‑booking in **one tap**. FAB sits in thumb reach (bottom‑right). |
| **Service Select** | • Icon + name + duration cards<br>• 2‑column grid | Tap area ≥ 48 px; highlight selected card. |
| **Date & Time** | • Horizontal 7‑day calendar<br>• Available time **chips** | Quick‑filter chips: “Earliest”, “Evening”, etc. |
| **Stylist** | • “No preference”, “Same as last” shortcuts<br>• Stylist cards + ★rating | Show extra fee on card if applicable. |
| **Confirm** | • Summary + total price<br>• Primary button **Reserve**<br>• Secondary **Edit** | Button fixed at **safe‑area bottom**; give feedback < 100 ms. |

---

## 3. Smoothing the Experience

1. **1‑tap Rebook** – swipe right on a past reservation to auto‑fill next slot  
2. **Skeleton loaders** – oval & text placeholders to cut perceived wait  
3. **Rich push notifications** – reminders (opt‑in) day before & morning of booking  
4. **Offline fallback** – store locally when offline, auto‑sync on reconnect  
5. **Accessibility** – semantic HTML first, focus ring visible, contrast ≥ 4.5 : 1  

---

## 4. Why This Feels “Snappy”

| Aspect | Concrete measures |
| ------ | ----------------- |
| **Fewest taps** | Service → Date → Stylist → Confirm – *max 4 taps* to finish |
| **Shortest thumb travel** | Bottom Sheets keep main actions in the thumb zone |
| **Lightweight load** | Code‑split non‑home views; JS budget \< 70 KB Gzip |
| **Network resilience** | Service & staff JSON pre‑cached via Service Worker |

---

## 5. Implementation Notes

| Item | Recommended stack |
| ---- | ----------------- |
| **UI Framework** | React with Preact Signals **or** SolidStart (easy Islands) |
| **Routing** | In‑page state machine (XState/Eigenstate); decouple from URL |
| **State** | Zustand/Recoil + IndexedDB → Background Sync API |
| **Performance budget** | JS \< 70 KB (gz), LCP \< 2.5 s, CLS \< 0.1 |

---

### Summary
> **“Home + 3 Bottom‑Sheet steps + Confirmation”** delivers the fastest, most finger‑friendly path from intent to booked appointment, even on slow networks or with one‑hand use.
