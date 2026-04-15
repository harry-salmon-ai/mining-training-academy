package seed

import (
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedTimeUtilizationModule(database *gorm.DB, force bool) {
	var moduleCount int64
	database.Model(&models.Module{}).Where("slug = ?", "time-utilization-framework").Count(&moduleCount)
	if moduleCount > 0 {
		if !force {
			log.Println("Time Utilization module already exists, skipping")
			return
		}
		forceDeleteModule(database, "time-utilization-framework")
	}

	log.Println("Seeding Time Utilization Framework module...")

	var category models.Category
	database.Where("slug = ?", "mining-operations").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed TUM: category or admin not found")
		return
	}

	desc := "The Time Utilization Framework (TUM) — how mining equipment spends every hour, from Calendar Time through Availability, Utilization, and Effective Utilization. Includes KPI calculations (MTBF, MTTR), the MTS perspective on modernising the TUM for operational technology, and practical implementation guidance."
	dur := 60
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Time Utilization Framework",
		Slug:        "time-utilization-framework",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelFoundation,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"TUM", "time usage model", "availability", "utilization", "KPI", "MTBF", "MTTR", "operational technology", "fleet management"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 1: Introduction ─────────────────────────────────────────
	s1 := createSection(database, module.ID, "Introduction", 0)

	createSlide(database, s1.ID, "Time Utilization Framework", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Time Utilization Framework",
		"subtitle": "Measuring how mining equipment spends every hour",
		"bullets": []string{
			"Why tracking equipment time is critical to mining profitability",
			"The GMG hierarchy from Calendar Time to Productive Time",
			"KPI calculations — Availability, Utilization, MTBF, MTTR, and beyond",
		},
	})

	createSlide(database, s1.ID, "Why Time Utilization Matters", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Why Time Utilization Matters",
		"introText": "Every minute of untracked downtime is a hidden cost. A single haul truck can cost $400–600 per operating hour — if 15% of its scheduled time disappears into poorly categorised delays, that is millions of dollars per year per truck that nobody is accountable for.",
		"items": []map[string]interface{}{
			{
				"title":       "Consistent KPIs",
				"description": "The TUM provides a standardised language for measuring performance. When every site defines 'availability' the same way, benchmarking becomes meaningful.",
			},
			{
				"title":       "Cross-Site Benchmarking",
				"description": "Mining companies with multiple operations need apples-to-apples comparisons. The TUM makes this possible by enforcing a common time classification hierarchy.",
			},
			{
				"title":       "Data-Driven Decisions",
				"description": "Without a TUM, improvement efforts are guesswork. With one, you can pinpoint exactly where time is being lost — downtime, standby, or operating delays — and direct resources accordingly.",
			},
		},
	})

	// ── Section 2: The Traditional TUM Hierarchy ────────────────────────
	s2 := createSection(database, module.ID, "The Traditional TUM Hierarchy", 1)

	createSlide(database, s2.ID, "Calendar Time Breakdown", "DIAGRAM", 0, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The TUM Hierarchy — Calendar Time Breakdown",
		"description": "The Global Mining Group (GMG) framework starts with Calendar Time and progressively subdivides it into smaller, more specific time blocks. Every hour in a period must land in exactly one leaf category — there are no gaps or overlaps.",
		"diagramType": "tree",
		"diagramData": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "CT", "label": "Calendar Time", "description": "Total hours in the period (e.g., 8,760 hrs/yr)"},
				{"id": "ST", "label": "Scheduled Time", "parent": "CT", "description": "Planned operating window"},
				{"id": "UST", "label": "Unscheduled Time", "parent": "CT", "description": "Holidays, shutdowns, no-production days"},
				{"id": "AT", "label": "Available Time", "parent": "ST", "description": "Equipment is functional and ready"},
				{"id": "DT", "label": "Down Time", "parent": "ST", "description": "Equipment required but not functional"},
				{"id": "OT", "label": "Operating Time", "parent": "AT", "description": "Equipment actively working"},
				{"id": "SB", "label": "Standby", "parent": "AT", "description": "Available but not operating"},
				{"id": "SDT", "label": "Scheduled Down", "parent": "DT", "description": "Planned maintenance, services"},
				{"id": "UDT", "label": "Unscheduled Down", "parent": "DT", "description": "Breakdowns, failures"},
				{"id": "PT", "label": "Productive Time", "parent": "OT", "description": "Hauling, loading, dumping"},
				{"id": "OD", "label": "Operating Delays", "parent": "OT", "description": "Fuelling, operator change, queuing"},
			},
		},
	})

	createSlide(database, s2.ID, "Downtime Categories", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Downtime Categories",
		"introText": "Downtime is when equipment is required for operation but cannot function. It splits into two fundamentally different buckets that drive very different management responses.",
		"items": []map[string]interface{}{
			{
				"title":       "Scheduled Downtime",
				"description": "Planned maintenance activities — services, component changeouts, inspections. These are known in advance and can be optimised around shift changes or low-demand periods. The goal is not to eliminate scheduled downtime but to make it efficient and predictable.",
			},
			{
				"title":       "Unscheduled Downtime",
				"description": "Breakdowns and failures that were not planned — engine faults, tyre blowouts, hydraulic leaks. This is the maintenance department's primary enemy. Every unscheduled event has a direct cost (repair) plus an indirect cost (lost production).",
			},
			{
				"title":       "The Maintenance Perspective",
				"description": "Maintenance teams are measured on downtime. Their goal is to shift the ratio from unscheduled (reactive) to scheduled (proactive). A world-class operation targets <5% unscheduled downtime as a proportion of total downtime.",
			},
		},
	})

	createSlide(database, s2.ID, "Standby Categories", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Standby Categories",
		"introText": "Standby is one of the most misunderstood categories. The machine IS available — it could run — but it is not operating. Standby is available time, NOT downtime. Confusing the two leads to incorrect KPIs and misallocated blame.",
		"items": []map[string]interface{}{
			{
				"title":       "Operational Standby",
				"description": "The machine is ready but cannot operate due to external factors — blast clearance, weather delays, no operator available, waiting for a loading unit. These are typically short-duration and operations-owned.",
			},
			{
				"title":       "Management / Surplus Standby",
				"description": "Equipment is available but not required — low demand, excess fleet, or deliberate park-up. This is a fleet-sizing and planning decision, not an equipment problem.",
			},
			{
				"title":       "Why It Matters",
				"description": "High standby with high availability means the operation has more fleet capacity than it is using. This is an operations and planning problem, not a maintenance problem. Misclassifying standby as downtime inflates maintenance KPIs unfairly.",
			},
		},
	})

	// ── Section 3: Availability ─────────────────────────────────────────
	s3 := createSection(database, module.ID, "Core KPIs — Availability", 2)

	createSlide(database, s3.ID, "Understanding Availability", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Understanding Availability",
		"introText": "Availability is the most widely reported TUM metric — but there are three distinct flavours, each answering a different question. Using the wrong one leads to misleading conclusions.",
		"items": []map[string]interface{}{
			{
				"title":       "Mechanical (Physical) Availability",
				"description": "Available Time ÷ (Available Time + Downtime). Answers: 'When we needed the machine, could it run?' This is the maintenance department's primary scorecard metric. It excludes unscheduled time entirely — it only measures the maintenance team's ability to keep equipment functional.",
			},
			{
				"title":       "Scheduled Availability",
				"description": "Available Time ÷ Scheduled Time. Factors in both downtime and the full scheduled window. Useful for comparing performance across different roster structures (e.g., 24/7 vs 2×12-hour shifts).",
			},
			{
				"title":       "Calendar Availability",
				"description": "Available Time ÷ Calendar Time. The broadest view — includes unscheduled time such as shutdowns, holidays, and no-production days. Used for fleet sizing and long-range mine planning.",
			},
		},
		"keyFact": "A truck can have 92% mechanical availability but only 78% calendar availability if the site runs a 5-day roster. Always check which availability is being quoted.",
	})

	createSlide(database, s3.ID, "What Drives Availability?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What Drives Availability?",
		"introText": "Availability is an output, not a lever. To improve it you need to work on the inputs — the frequency and duration of downtime events.",
		"items": []map[string]interface{}{
			{
				"title":       "Reduce Unscheduled Breakdowns",
				"description": "Reliability improvement through condition monitoring, oil analysis, vibration monitoring, and component life management. Preventing failures is always cheaper than fixing them.",
			},
			{
				"title":       "Shorten Repair Times",
				"description": "Parts availability, technician skill level, workshop layout, and tooling all affect MTTR. Every hour saved in the workshop is an hour returned to production.",
			},
			{
				"title":       "Optimise Scheduled Maintenance",
				"description": "Perform services during shift changes or low-demand windows. Extend component life where safely possible. Use reliability-centred maintenance (RCM) to challenge unnecessary PM tasks.",
			},
			{
				"title":       "Common Pitfall",
				"description": "Chasing availability percentage without checking whether the extra hours actually get utilised. A fleet at 95% availability but 70% utilization has an operations problem, not a maintenance win.",
			},
		},
	})

	// ── Section 4: Utilization & Effective Utilization ───────────────────
	s4 := createSection(database, module.ID, "Core KPIs — Utilization & Effective Utilization", 3)

	createSlide(database, s4.ID, "Understanding Utilization", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Understanding Utilization",
		"introText": "Utilization measures how well available time is converted into actual operating time. It is the operations team's metric — not maintenance's.",
		"items": []map[string]interface{}{
			{
				"title":       "Utilization Formula",
				"description": "Operating Time ÷ Available Time. Answers: 'Of the time the machine could run, how much did it actually work?' Everything that reduces utilization is a standby or delay event that happened while the machine was available.",
			},
			{
				"title":       "Low Utilization + High Availability",
				"description": "This is an operations problem — excess fleet, poor dispatch allocation, too much standby. The machines are ready but the operation is not using them. Look at fleet sizing, dispatch rules, and shift handover delays.",
			},
			{
				"title":       "High Utilization + Low Availability",
				"description": "This is a maintenance problem — not enough uptime, but every available hour is being consumed. The operation is starved for capacity. Priority is improving reliability and reducing repair times.",
			},
		},
	})

	createSlide(database, s4.ID, "Effective Utilization (Use of Availability)", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Effective Utilization — Use of Availability",
		"introText": "Effective Utilization is the single best indicator of how productively a fleet's scheduled time is being spent. It forces both maintenance AND operations to own the result jointly.",
		"items": []map[string]interface{}{
			{
				"title":       "The Formula",
				"description": "Effective Utilization = Operating Time ÷ Scheduled Time. This is mathematically equivalent to Availability × Utilization. Sometimes called 'Use of Availability' or 'Overall Time Efficiency'.",
			},
			{
				"title":       "Worked Example",
				"description": "Availability = 90%, Utilization = 85%. Effective Utilization = 0.90 × 0.85 = 76.5%. That means nearly a quarter of scheduled time is lost to either downtime or standby. Neither department can point fingers — both contribute to the gap.",
			},
			{
				"title":       "Why It Matters",
				"description": "Reporting availability and utilization separately lets each department optimise in isolation. Effective Utilization makes the combined result visible and creates shared accountability for the total productive output of the fleet.",
			},
		},
		"keyFact": "If Availability = 90% and Utilization = 85%, Effective Utilization = 76.5% — nearly a quarter of scheduled time is lost.",
	})

	createSlide(database, s4.ID, "KPI Relationship Map", "DIAGRAM", 2, map[string]interface{}{
		"type":        "diagram",
		"heading":     "KPI Relationship Map",
		"description": "Every KPI maps to a specific ratio in the TUM hierarchy. The key insight is the multiplication effect — Availability × Utilization = Effective Utilization. Each department owns a different segment of the tree.",
		"diagramType": "kpi-map",
		"diagramData": map[string]interface{}{
			"ratios": []map[string]interface{}{
				{
					"name":        "Mechanical Availability",
					"formula":     "Available Time ÷ (Available Time + Downtime)",
					"owner":       "Maintenance",
					"description": "Can the machine run when needed?",
				},
				{
					"name":        "Utilization",
					"formula":     "Operating Time ÷ Available Time",
					"owner":       "Operations / Dispatch",
					"description": "Of available time, how much is used?",
				},
				{
					"name":        "Effective Utilization",
					"formula":     "Operating Time ÷ Scheduled Time = Availability × Utilization",
					"owner":       "Joint — Maintenance + Operations",
					"description": "How productively is scheduled time spent?",
				},
				{
					"name":        "Calendar Availability",
					"formula":     "Available Time ÷ Calendar Time",
					"owner":       "Planning / Management",
					"description": "Fleet capacity across the full calendar",
				},
			},
			"callout": map[string]interface{}{
				"text":        "Availability × Utilization = Effective Utilization",
				"description": "The multiplication effect means a 5% improvement in either metric has a compounding impact on overall fleet productivity.",
			},
		},
	})

	// ── Section 5: Reliability KPIs ─────────────────────────────────────
	s5 := createSection(database, module.ID, "Reliability KPIs — MTBF & MTTR", 4)

	createSlide(database, s5.ID, "MTBF & MTTR", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "MTBF & MTTR — Reliability and Maintainability",
		"introText": "While availability tells you the result, MTBF and MTTR tell you why. They are the diagnostic metrics that reliability and maintenance engineers use to drive improvement.",
		"items": []map[string]interface{}{
			{
				"title":       "MTBF — Mean Time Between Failures",
				"description": "Operating Hours ÷ Number of Failures. Higher is better. Measures how reliable the equipment is — how long it runs before something breaks. Improving MTBF means preventing failures through condition monitoring, proactive maintenance, and design improvements.",
			},
			{
				"title":       "MTTR — Mean Time To Repair",
				"description": "Total Repair Hours ÷ Number of Repairs. Lower is better. Measures how quickly failures are resolved — parts availability, technician skill, diagnostic tools, and workshop efficiency all feed into MTTR.",
			},
			{
				"title":       "The Availability Equation",
				"description": "Availability = MTBF ÷ (MTBF + MTTR). This shows that improving either metric improves availability, but MTBF improvements (preventing failures) have a bigger long-term payoff than faster repairs. A machine that never breaks down doesn't need fast repairs.",
			},
			{
				"title":       "Reactive → Proactive",
				"description": "Low MTBF with low MTTR means the team is good at fixing things but things break too often (reactive). High MTBF with moderate MTTR means failures are rare and repairs are manageable (proactive). World-class operations target high MTBF as the primary goal.",
			},
		},
	})

	// ── Section 5b: Knowledge Check ─────────────────────────────────────
	s5b := createSection(database, module.ID, "Knowledge Check", 5)

	createSlide(database, s5b.ID, "Classify These Events", "INTERACTIVE", 0, map[string]interface{}{
		"type":    "interactive",
		"heading": "Knowledge Check — Classify These Events",
		"interactiveType": "classify",
		"instructions":    "For each scenario, decide whether the event should be classified as Downtime, Standby, or Operating Delay in the TUM.",
		"scenarios": []map[string]interface{}{
			{
				"event":       "Blast clearance — all trucks parked for 45 minutes",
				"answer":      "Standby",
				"explanation": "The trucks are available (not broken) but cannot operate due to an external event. This is operational standby.",
			},
			{
				"event":       "Flat tyre on a haul truck — 3 hours to change",
				"answer":      "Downtime (Unscheduled)",
				"explanation": "The truck cannot operate due to an equipment failure. This is unscheduled downtime — a breakdown.",
			},
			{
				"event":       "No operator available at shift change — truck idle for 30 minutes",
				"answer":      "Standby",
				"explanation": "The truck is functional but there is no one to drive it. This is operational standby, not downtime.",
			},
			{
				"event":       "Scheduled 500-hour service",
				"answer":      "Downtime (Scheduled)",
				"explanation": "Planned maintenance. The truck is required but out of service for a known, scheduled activity.",
			},
			{
				"event":       "Fleet Management System offline — truck running but not dispatched",
				"answer":      "Standby or Tech Downtime (debated)",
				"explanation": "This is the grey area the MTS article highlights. Traditionally classified as standby, but increasingly argued to be a technology-related downtime event.",
			},
		},
		"calculationQuestion": map[string]interface{}{
			"question": "Availability is 88%, Utilization is 82%. What is Effective Utilization?",
			"answer":   "72.2%",
			"formula":  "0.88 × 0.82 = 0.7216 = 72.2%",
		},
	})

	// ── Section 6: The Technology Problem ────────────────────────────────
	s6 := createSection(database, module.ID, "The Technology Problem — MTS Perspective", 6)

	createSlide(database, s6.ID, "The Tech Revolution and the TUM", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "The Tech Revolution and the TUM's Ticking Clock",
		"introText": "As mining becomes more technology-dependent, the traditional TUM is being challenged. Autonomous Haulage Systems (AHS), Fleet Management Systems (FMS), Collision Avoidance (CAS), and Fatigue Detection are not part of a machine's core mechanical, hydraulic, or electrical function — yet their failure can bring operations to a halt.",
		"items": []map[string]interface{}{
			{
				"title":       "The Gap in the Traditional TUM",
				"description": "These Operational Technology (OT) systems have no clear home in the traditional time classification. When a CAS system locks out a truck, is it downtime? Standby? The TUM doesn't say — and different sites classify it differently, destroying benchmarking consistency.",
			},
			{
				"title":       "Unclear Ownership",
				"description": "Who is responsible for technology uptime? It often falls between departments — maintenance doesn't own the software, IT doesn't own the on-machine hardware, and operations just sees a truck that can't work. Few metrics exist and fewer people are accountable.",
			},
			{
				"title":       "The MTS Insight",
				"description": "Mine Tech Services (MTS) identified that the allocation of Operational Technology downtime is 'hotly debated' across the industry. A discussion with mining technology professionals revealed fundamentally different views on how to classify these events — highlighting the urgent need for an updated framework.",
			},
		},
		"keyFact": "Technology failures can halt production for safety, productivity, or regulatory reasons — but the traditional TUM has no category for them.",
	})

	createSlide(database, s6.ID, "The Great Debate: Where Does Tech Downtime Go?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Great Debate — Where Does Tech Downtime Go?",
		"introText": "MTS highlighted three competing views from industry professionals. Each has merits and trade-offs.",
		"items": []map[string]interface{}{
			{
				"title":       "View 1: Classify as True Downtime",
				"description": "If a system is required by site regulations or company standards for effective production, its failure is a functional failure and the equipment should be classified as down. This gives urgency to resolve the issue — just like any mechanical breakdown. Risk: it penalises the maintenance department for issues outside their control.",
			},
			{
				"title":       "View 2: 'Other Down Time' with Dedicated Codes",
				"description": "Contain technology failures within an 'Other Down Time' category with dedicated codes for FMS, HP-GPS, Wi-Fi, CAS, etc. This allows separate tracking without inflating the mechanical downtime KPIs. Risk: if it's not in the main downtime number, it may not get management attention.",
			},
			{
				"title":       "View 3: Create an Entirely New Category",
				"description": "Add a new branch to the TUM — 'Downtime - Technology' — sitting alongside Scheduled and Unscheduled mechanical downtime. This gives technology its own metrics, ownership, and visibility without distorting existing KPIs. This is the direction the GMG is moving.",
			},
		},
		"keyFact": "The key principle: technology stoppages should NOT dilute the wrong department's metrics. If maintenance doesn't own the tech, tech downtime should not reduce their availability score.",
	})

	// ── Section 7: A Modernised TUM ─────────────────────────────────────
	s7 := createSection(database, module.ID, "A Modernised TUM", 7)

	createSlide(database, s7.ID, "Updated TUM with Technology Branch", "DIAGRAM", 0, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Modernised TUM — Technology Downtime Branch",
		"description": "The updated framework adds a 'Downtime - Technology' branch alongside the traditional Scheduled and Unscheduled mechanical downtime categories. This preserves existing KPIs while creating visibility and accountability for technology-related production losses.",
		"diagramType": "tree",
		"diagramData": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "ST", "label": "Scheduled Time", "description": "Planned operating window"},
				{"id": "AT", "label": "Available Time", "parent": "ST"},
				{"id": "DT", "label": "Down Time", "parent": "ST"},
				{"id": "SDT", "label": "Scheduled Down", "parent": "DT", "description": "Planned maintenance"},
				{"id": "UDT", "label": "Unscheduled Down", "parent": "DT", "description": "Breakdowns"},
				{"id": "TDT", "label": "Technology Down", "parent": "DT", "description": "FMS, CAS, AHS, GPS, Wi-Fi failures", "highlight": true},
			},
			"callout": map[string]interface{}{
				"text":        "New: Downtime — Technology",
				"description": "Separate tracking with dedicated codes, ownership, and KPIs. Excluded from mechanical availability calculations.",
			},
		},
	})

	createSlide(database, s7.ID, "New Technology KPIs", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "New Technology KPIs — MTS Recommendations",
		"introText": "With a dedicated technology downtime category, new KPIs become possible. MTS recommends three core metrics to drive accountability and continuous improvement for onboard technology systems.",
		"items": []map[string]interface{}{
			{
				"title":       "MTBS-Tech (Mean Time Between Technology Stoppages)",
				"description": "Operating Hours ÷ Technology-Related Stoppage Count. The technology equivalent of MTBF. Measures how often tech systems cause a production interruption. Higher is better.",
			},
			{
				"title":       "Technology Downtime Hours",
				"description": "Sum of all downtime hours where the delay code is technology-related. The raw volume metric — how much total production time is lost to tech failures per period. Trend this monthly to justify investment in technology maintenance resources.",
			},
			{
				"title":       "Time Between Tech PM Inspections",
				"description": "Tracks the interval between preventive maintenance inspections of onboard technology. MTS emphasises this metric drives a proactive approach to technology upkeep — moving from 'fix it when it breaks' to regular inspection cycles, just like mechanical PMs.",
			},
		},
		"keyFact": "Technology-related stoppages should be excluded from the maintenance department's availability calculations. If maintenance doesn't maintain the tech, the tech bucket should not dilute their metrics.",
	})

	// ── Section 8: Practical Application ────────────────────────────────
	s8 := createSection(database, module.ID, "Practical Application", 8)

	createSlide(database, s8.ID, "Implementing TUM at Your Site", "PROCESS", 0, map[string]interface{}{
		"type":        "process",
		"heading":     "Implementing TUM at Your Site",
		"description": "Rolling out or updating a Time Usage Model is a cross-departmental effort. Technology, operations, maintenance, and planning must align on definitions before any data is meaningful.",
		"steps": []map[string]interface{}{
			{
				"label":       "Step 1 — Define Time Codes",
				"description": "Bring all departments together to agree on time categories, delay codes, and ownership. Include a dedicated technology downtime category with specific codes for FMS, CAS, GPS, Wi-Fi, and AHS.",
			},
			{
				"label":       "Step 2 — Configure in FMS / Dispatch",
				"description": "Programme the agreed codes into the Fleet Management System so operators can select them in real time. Ensure the system enforces a code selection — no unclassified time allowed.",
			},
			{
				"label":       "Step 3 — Train Operators",
				"description": "Operators are the front line of data quality. Train them on what each code means and when to use it. Poor code selection at the operator level corrupts every KPI downstream.",
			},
			{
				"label":       "Step 4 — Build Dashboards",
				"description": "Create real-time and shift-summary dashboards showing availability, utilization, effective utilization, and technology downtime. Make the data visible to everyone — from the control room to the boardroom.",
			},
			{
				"label":       "Step 5 — Review and Iterate Monthly",
				"description": "Hold monthly reviews with all stakeholders. Challenge misclassified events, identify emerging patterns, and adjust codes or thresholds as the operation evolves.",
			},
		},
	})

	createSlide(database, s8.ID, "Traditional vs Modern TUM", "COMPARISON", 1, map[string]interface{}{
		"type":    "comparison",
		"heading": "Traditional TUM vs Modern TUM",
		"items": []map[string]interface{}{
			{
				"name":        "Traditional TUM",
				"description": "The GMG framework as originally published — suited to mechanical fleets with no onboard technology dependency.",
				"pros": []string{
					"Universally understood across the industry",
					"Calendar → Scheduled → Available → Operating hierarchy",
					"Core KPIs: Availability, Utilization, MTBF, MTTR",
					"Sufficient for purely mechanical fleets",
				},
				"cons": []string{
					"No category for technology-related stoppages",
					"Technology failures lumped into general downtime or standby",
					"No dedicated ownership or metrics for OT systems",
					"Benchmarking is distorted when sites classify tech events differently",
				},
			},
			{
				"name":        "Modern TUM (with Technology Branch)",
				"description": "The updated framework adds a dedicated Technology Downtime branch alongside Scheduled and Unscheduled maintenance — required for any site running operational technology.",
				"pros": []string{
					"Same core hierarchy preserved — no KPI disruption",
					"Downtime adds a third branch: Technology Down",
					"Dedicated delay codes for FMS, CAS, AHS, GPS, Wi-Fi",
					"New KPIs: MTBS-Tech, Tech Downtime Hours, Tech PM Interval",
					"Clear ownership: tech team owns tech metrics, maintenance owns mechanical",
				},
				"cons": []string{
					"Requires cross-department agreement on code definitions",
					"FMS configuration needed for new delay codes",
					"Not yet an official GMG standard — industry adoption is in progress",
				},
			},
		},
	})

	// ── Section 9: Assessment ───────────────────────────────────────────
	s9 := createSection(database, module.ID, "Assessment", 9)

	quizDesc := "Test your understanding of the Time Utilization Framework, KPI calculations, and the modernised TUM for operational technology."
	quiz := models.Quiz{
		ID:              uuid.New().String(),
		Title:           "Time Utilization Framework Assessment",
		Description:     &quizDesc,
		PassMark:        80,
		MaxAttempts:     3,
		ShuffleQuestions: true,
		ShuffleOptions:  true,
		ShowResults:     true,
		SectionID:       &s9.ID,
		ModuleID:        &module.ID,
	}
	database.Create(&quiz)

	questions := []struct {
		question    string
		qType       models.QuestionType
		explanation string
		options     []struct {
			text      string
			isCorrect bool
		}
	}{
		{
			question:    "What is Calendar Time in the TUM?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Calendar Time is the total hours in the period — e.g., 8,760 hours in a year, 720 hours in a 30-day month. It is the starting point of the entire hierarchy.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Total hours the equipment was operating", false},
				{"Total hours in the period (e.g., 8,760 hrs/year)", true},
				{"Total hours the equipment was available", false},
				{"Total scheduled production hours", false},
			},
		},
		{
			question:    "What is the formula for Mechanical (Physical) Availability?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Mechanical Availability = Available Time ÷ (Available Time + Downtime). It measures the maintenance team's ability to keep equipment functional when it is needed.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Operating Time ÷ Available Time", false},
				{"Available Time ÷ Calendar Time", false},
				{"Available Time ÷ (Available Time + Downtime)", true},
				{"Operating Time ÷ Scheduled Time", false},
			},
		},
		{
			question:    "A haul truck is parked during a blast clearance for 45 minutes. How should this be classified?",
			qType:       models.QuestionMultipleChoice,
			explanation: "The truck is available (not broken) but cannot operate due to an external event. This is operational standby — available time that is not operating time.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Unscheduled Downtime", false},
				{"Operating Delay", false},
				{"Standby", true},
				{"Scheduled Downtime", false},
			},
		},
		{
			question:    "What is the Utilization formula, and which department owns it?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Utilization = Operating Time ÷ Available Time. It measures how well the operations team converts available hours into working hours. Maintenance owns availability; operations owns utilization.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Available Time ÷ Scheduled Time — owned by Maintenance", false},
				{"Operating Time ÷ Available Time — owned by Operations", true},
				{"Operating Time ÷ Calendar Time — owned by Planning", false},
				{"Operating Time ÷ Scheduled Time — owned by Maintenance", false},
			},
		},
		{
			question:    "If Availability is 88% and Utilization is 82%, what is Effective Utilization?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Effective Utilization = Availability × Utilization = 0.88 × 0.82 = 0.7216 = 72.2%.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"85.0%", false},
				{"72.2%", true},
				{"78.5%", false},
				{"66.0%", false},
			},
		},
		{
			question:    "What does MTBF measure, and is higher or lower better?",
			qType:       models.QuestionMultipleChoice,
			explanation: "MTBF = Mean Time Between Failures = Operating Hours ÷ Number of Failures. Higher is better — it means the equipment runs longer between breakdowns, indicating better reliability.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Mean repair time — lower is better", false},
				{"Mean time between failures — higher is better", true},
				{"Mean time between failures — lower is better", false},
				{"Mean operating hours — lower is better", false},
			},
		},
		{
			question:    "Availability = MTBF ÷ (MTBF + MTTR). Which improvement has a bigger long-term payoff?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Improving MTBF (preventing failures from occurring) has a bigger long-term payoff than improving MTTR (fixing failures faster). A machine that never breaks down doesn't need fast repairs.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Improving MTTR (faster repairs)", false},
				{"Improving MTBF (preventing failures)", true},
				{"Both have equal impact", false},
				{"Neither — availability is fixed by fleet size", false},
			},
		},
		{
			question:    "According to MTS, why shouldn't technology failures be lumped into general downtime?",
			qType:       models.QuestionMultipleChoice,
			explanation: "MTS argues that if maintenance doesn't own the technology, tech-related stoppages should not dilute their availability metrics. Lumping them together misrepresents mechanical performance and removes accountability from the team that should own the tech.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Technology failures are too rare to track", false},
				{"It dilutes the wrong department's metrics", true},
				{"General downtime codes are more accurate", false},
				{"Technology systems are not safety-critical", false},
			},
		},
		{
			question:    "What is MTBS-Tech?",
			qType:       models.QuestionMultipleChoice,
			explanation: "MTBS-Tech = Mean Time Between Technology-Related Stoppages = Operating Hours ÷ Technology Stoppage Count. It is the technology equivalent of MTBF.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Mean Time Between Scheduled maintenance", false},
				{"Mean Time Between Technology-Related Stoppages", true},
				{"Maximum Time Before System failure", false},
				{"Mean Total Breakdown Severity — Technology", false},
			},
		},
		{
			question:    "A truck has 90% mechanical availability, but the site runs a 5-day roster. Will its calendar availability be higher or lower than 90%?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Calendar availability includes unscheduled time (weekends, in this case). Available Time stays the same but Calendar Time is larger than Scheduled Time, so the percentage drops. Calendar Availability will be lower than mechanical availability.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Higher — more total hours means higher availability", false},
				{"Lower — calendar time is larger so the percentage drops", true},
				{"The same — roster doesn't affect availability", false},
				{"It depends on the number of breakdowns", false},
			},
		},
	}

	for i, q := range questions {
		question := models.QuizQuestion{
			ID:        uuid.New().String(),
			QuizID:    quiz.ID,
			Type:      q.qType,
			Question:  q.question,
			Points:    1,
			SortOrder: i,
		}
		if q.explanation != "" {
			question.Explanation = &q.explanation
		}
		database.Create(&question)

		for j, o := range q.options {
			database.Create(&models.QuizOption{
				ID:         uuid.New().String(),
				QuestionID: question.ID,
				Text:       o.text,
				IsCorrect:  o.isCorrect,
				SortOrder:  j,
			})
		}
	}

	createSlide(database, s9.ID, "TUM Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Time Utilization Framework Assessment",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{
				"question":     "What is Mechanical Availability?",
				"options":      []string{"Operating Time ÷ Available Time", "Available Time ÷ Calendar Time", "Available Time ÷ (Available + Downtime)", "Operating Time ÷ Scheduled Time"},
				"correctIndex": 2,
				"explanation":  "Mechanical Availability = Available Time ÷ (Available Time + Downtime).",
			},
			{
				"question":     "If Availability = 88% and Utilization = 82%, Effective Utilization is:",
				"options":      []string{"85.0%", "72.2%", "78.5%", "66.0%"},
				"correctIndex": 1,
				"explanation":  "0.88 × 0.82 = 72.2%.",
			},
			{
				"question":     "Why shouldn't tech failures go in general downtime?",
				"options":      []string{"Too rare to track", "Dilutes the wrong department's metrics", "General codes are more accurate", "Tech isn't safety-critical"},
				"correctIndex": 1,
				"explanation":  "MTS: if maintenance doesn't own the tech, tech stoppages should not dilute their availability score.",
			},
		},
	})

	// ── Section 10: Completion ──────────────────────────────────────────
	s10 := createSection(database, module.ID, "Completion", 10)

	createSlide(database, s10.ID, "Module Complete", "COMPLETION", 0, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "Congratulations — you have completed the Time Utilization Framework module.",
		"keyTakeaways": []string{
			"The TUM breaks Calendar Time into a hierarchy: Scheduled → Available → Operating, with Downtime and Standby as the key loss categories",
			"Mechanical Availability = Available Time ÷ (Available + Downtime) — the maintenance department's primary scorecard metric",
			"Utilization = Operating Time ÷ Available Time — owned by operations, not maintenance",
			"Effective Utilization = Availability × Utilization — the single best indicator of fleet productivity, creating shared accountability",
			"MTBF and MTTR are diagnostic metrics: Availability = MTBF ÷ (MTBF + MTTR) — preventing failures beats faster repairs",
			"The modern TUM must include a Technology Downtime branch with dedicated KPIs (MTBS-Tech) to capture operational technology failures without diluting mechanical metrics",
		},
		"showCertificate": true,
	})

	log.Println("Time Utilization module seeded with 10 sections, 19 slides, and 10-question quiz")
}
