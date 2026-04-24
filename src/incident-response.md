---
title: "Incident Response"
weight: 30
---

# Incident Response Procedure

**Owner:** {{it_lead}}
**Applies to:** All staff at {{company_name}}
**Last Reviewed:** {{last_reviewed}}

---

## Purpose

This procedure defines how {{company_name}} identifies, responds to, and recovers from IT and security incidents to minimise impact and prevent recurrence.

---

## Scope

This procedure covers all IT and security incidents affecting {{company_name}} systems, data, or operations — including but not limited to data breaches, ransomware, service outages, and physical security events.

---

## Severity Levels

| Level | Name | Description | Response Time |
|-------|------|-------------|--------------|
| P1 | Critical | Full outage or confirmed data breach | Immediate (< 1 hour) |
| P2 | High | Partial outage, suspected breach, significant data loss | < 4 hours |
| P3 | Medium | Degraded service, isolated issue, no data loss | < 24 hours |
| P4 | Low | Minor issue, workaround available | Next business day |

---

## Incident Response Team

| Role | Person | Contact |
|------|--------|---------|
| Incident Commander | {{it_lead}} | Via {{company_short}} internal comms |
| Compliance Lead | {{compliance_officer}} | {{compliance_officer_email}} |
| Executive Sponsor | {{ceo_name}} | Via {{company_short}} internal comms |

---

## Procedure

### Phase 1 — Identification

1. Any staff member who suspects an incident must report it immediately to {{it_lead}}
2. {{it_lead}} logs the incident in the Incident Register with:
   - Date/time identified
   - Systems affected
   - Initial severity assessment
   - Reporter name

### Phase 2 — Containment

**Immediate actions ({{it_lead}}):**
1. Isolate affected systems from the network if required
2. Preserve evidence — do not power off systems without IT guidance
3. Revoke compromised credentials
4. Notify {{compliance_officer}} if personal data may be involved

**Communication:**
- P1/P2: Notify {{ceo_name}} within 1 hour
- P1 data breach: {{compliance_officer}} assesses regulatory notification requirements (72-hour window)

### Phase 3 — Investigation

1. {{it_lead}} leads root cause analysis
2. Document timeline of events
3. Identify scope of impact (systems, data, users affected)
4. Engage external forensics if required (approval from {{ceo_name}})

### Phase 4 — Recovery

1. Restore systems from clean backups where necessary
2. Apply patches or configuration fixes that address root cause
3. Reset all potentially compromised credentials
4. Verify restoration before returning systems to production
5. Confirm with {{compliance_officer}} before closing

### Phase 5 — Post-Incident Review

Within 5 business days of incident closure:
1. {{it_lead}} produces Post-Incident Report covering:
   - Timeline
   - Root cause
   - Impact summary
   - Actions taken
   - Recommendations to prevent recurrence
2. Review meeting with {{compliance_officer}} and {{ceo_name}}
3. Update this procedure if gaps identified

---

## Reporting a Suspected Incident

**Any staff member** can and should report a suspected incident:

- Email: {{compliance_officer_email}}
- Or contact {{it_lead}} directly via internal comms

Do not attempt to resolve or investigate an incident yourself. Report it and let the Incident Response Team take over.

---

## Related Documents

- Data Protection & Security Policy
- Business Continuity Plan
- Data Subject Rights Procedure

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | {{last_reviewed}} | {{doc_owner}} | Initial version |
