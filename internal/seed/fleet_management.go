package seed

import (
	"encoding/json"
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedFleetManagementModule(database *gorm.DB, force bool) {
	var existing models.Module
	database.Where("slug = ?", "fleet-management-systems").First(&existing)
	if existing.ID != "" {
		if force {
			log.Println("SEED=true: deleting existing Fleet Management Systems module for re-seed...")
			// Delete in dependency order
			database.Where("module_id = ?", existing.ID).Delete(&models.ModuleTag{})
			var sections []models.Section
			database.Where("module_id = ?", existing.ID).Find(&sections)
			for _, sec := range sections {
				var slides []models.Slide
				database.Where("section_id = ?", sec.ID).Find(&slides)
				for _, slide := range slides {
					database.Where("slide_id = ?", slide.ID).Delete(&models.SlideProgress{})
					database.Delete(&slide)
				}
			}
			var quizzes []models.Quiz
			database.Where("module_id = ?", existing.ID).Find(&quizzes)
			for _, quiz := range quizzes {
				var questions []models.QuizQuestion
				database.Where("quiz_id = ?", quiz.ID).Find(&questions)
				for _, q := range questions {
					database.Where("question_id = ?", q.ID).Delete(&models.QuizOption{})
					database.Delete(&q)
				}
				var attempts []models.QuizAttempt
				database.Where("quiz_id = ?", quiz.ID).Find(&attempts)
				for _, a := range attempts {
					database.Where("attempt_id = ?", a.ID).Delete(&models.QuizAnswer{})
					database.Delete(&a)
				}
				database.Delete(&quiz)
			}
			database.Where("module_id = ?", existing.ID).Delete(&models.Enrollment{})
			database.Where("module_id = ?", existing.ID).Delete(&models.Certificate{})
			database.Where("module_id = ?", existing.ID).Delete(&models.Section{})
			database.Delete(&existing)
			log.Println("Existing Fleet Management Systems module deleted")
		} else {
			log.Println("Fleet Management Systems module already exists, skipping")
			return
		}
	}

	log.Println("Seeding Fleet Management Systems module...")

	var category models.Category
	database.Where("slug = ?", "mining-operations").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed FMS: category or admin not found")
		return
	}

	desc := "Mining Fleet Management Systems (FMS) — how FMS functions as a platform to measure, manage, and optimise mine site operations. Covers system requirements, production records (including drill data and MWD integration), fleet visibility, dispatch optimisation, grade control and CMMS integration, and implementation. Designed for engineers building a rigorous foundation in FMS."
	dur := 90
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Mining Fleet Management Systems",
		Slug:        "fleet-management-systems",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelIntermediate,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"FMS", "fleet management", "dispatch", "production records", "payload", "optimisation", "KPI", "grade control", "CMMS", "measure manage optimise"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 1: Introduction ─────────────────────────────────────────────
	s1 := createSection(database, module.ID, "Introduction to Fleet Management", 0)

	fmsJSON, _ := json.Marshal(map[string]interface{}{
		"type":     "title",
		"heading":  "Mining Fleet Management Systems",
		"subtitle": "A platform to measure, manage, and optimise mine site operations",
		"bullets": []string{
			"What FMS is and why it matters across all mine types",
			"The value ladder: how each capability layer compounds productivity",
			"Where FMS sits in the mining value chain",
		},
	})
	database.Create(&models.Slide{ID: uuid.New().String(), Title: "Mining Fleet Management Systems", Type: "TITLE", Content: fmsJSON, SortOrder: 0, SectionID: s1.ID})

	createSlide(database, s1.ID, "What is a Fleet Management System?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What is a Fleet Management System?",
		"introText": "An FMS is the operational nervous system of a mine. Its fundamental purpose is to measure operational performance and provide the tools to manage a site's daily operations. At its core, every FMS — regardless of vendor or scale — does three things.",
		"items": []map[string]interface{}{
			{
				"title":       "Measure",
				"description": "Turn raw positional and sensor data into accurate production measurement: cycles, payloads, delays, fuel consumption, and equipment hours. Every shift, every truck, every load. Recording is the mechanism — measurement is the outcome. This data is the basis for all operational accountability and mine planning feedback.",
			},
			{
				"title":       "Manage",
				"description": "Give supervisors and planners the tools to act on that measurement in real time — route trucks, manage delays, enforce standards, and make decisions based on what is actually happening rather than what was planned.",
			},
			{
				"title":       "Optimise: Dynamic Assignment",
				"description": "Algorithmically assigning trucks to the optimal loader in real time (dispatch optimisation) is an optimise-level FMS capability — the domain of systems like Modular Mining DISPATCH and MineStar Fleet. It delivers the highest productivity gains but requires the Measure and Manage foundations to be solid first.",
			},
		},
	})

	createSlide(database, s1.ID, "The FMS Value Ladder", "PROCESS", 2, map[string]interface{}{
		"type":        "process",
		"heading":     "The FMS Value Ladder",
		"description": "A full-scale FMS delivers up to 30% fleet productivity improvement — but not all at once. Each layer compounds on the last. Most operations leave the top two layers on the table.",
		"steps": []map[string]interface{}{
			{
				"title":       "Step 1 — Track: Equipment Visibility",
				"duration":    "+5–8%",
				"description": "Simply knowing where machines are and counting cycles makes the worst inefficiencies visible — lost trucks, unquantified idle time, inaccurate shift records. Sites moving from radio dispatch to GPS tracking see immediate productivity gains without changing a single process.",
			},
			{
				"title":       "Step 2 — Measure: Accurate Production Records",
				"duration":    "+5–8%",
				"description": "Reliable cycle data, payload records, and delay classification enable data-driven management. Payload compliance tightens, chronic delay patterns become visible, and shift-to-shift accountability improves. You cannot manage what you cannot measure.",
			},
			{
				"title":       "Step 3 — Manage: Real-Time Operations",
				"duration":    "+8–10%",
				"description": "Supervisors acting on live data — reducing queue times, enforcing standards, managing shift handovers, and responding to delays as they happen rather than reviewing them the next morning. Fuel reduction (3–8%) and shift change efficiency gains land here.",
			},
			{
				"title":       "Step 4 — Optimise: Dynamic Dispatch",
				"duration":    "+5–8%",
				"tier":        "React",
				"description": "Algorithmic truck-to-loader assignment, continuous re-optimisation, and match factor management. Advanced capability — only extracted on the foundation of Steps 1–3. A dispatch system on poor data produces confidently wrong assignments.",
				"warnings": []string{
					"Advanced capability: requires Steps 1–3 embedded before activation",
				},
			},
		},
	})

	createSlide(database, s1.ID, "FMS in the Mining Value Chain", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "FMS in the Mining Value Chain",
		"description": "FMS is not just a load and haul tool. It operates across every stage where mobile equipment works — from the drill pattern to the mill feed. The value chain below shows where FMS is active, and what it contributes at each stage.",
		"diagramType": "process",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{
					"label":       "Exploration & Planning",
					"description": "Geological data, resource modelling, mine scheduling",
					"highlight":   false,
					"fmsRole":     "",
				},
				{
					"label":       "Drill & Blast",
					"description": "Drill pattern execution, fragmentation",
					"highlight":   true,
					"fmsRole":     "FMS tracks drill position, metres drilled, bit changes, and production against pattern targets",
				},
				{
					"label":       "Load & Haul",
					"description": "Excavation, truck haulage, material routing",
					"highlight":   true,
					"fmsRole":     "PRIMARY DOMAIN — cycle tracking, payload, dispatch, delay management, fleet utilisation",
				},
				{
					"label":       "Processing / Mill",
					"description": "Crushing, grinding, extraction",
					"highlight":   true,
					"fmsRole":     "FMS tracks tonnes delivered to ROM and crusher, timing of deliveries to support mill scheduling — mill feed rate itself is controlled by the process control system (DCS/SCADA), not FMS",
				},
				{
					"label":       "Product & Rehab",
					"description": "Final product, site rehabilitation",
					"highlight":   false,
					"fmsRole":     "",
				},
			},
			"callout": map[string]interface{}{
				"text":        "FMS touches every stage where mobile equipment operates",
				"description": "From the drill rig to the crusher feed — if a machine moves, FMS should know about it.",
			},
		},
	})

	createSlide(database, s1.ID, "Core vs Advanced FMS Capabilities", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Core vs Advanced FMS Capabilities",
		"introText": "Not all FMS deployments are equal. Understanding the distinction between core and advanced capabilities is critical when evaluating systems, setting implementation expectations, and benchmarking performance.",
		"items": []map[string]interface{}{
			{
				"title":       "Core (any FMS, any operation)",
				"description": "GPS tracking, activity state inference, cycle recording, shift production reports, delay logging, basic utilisation metrics, onboard operator display, fuel consumption recording. Also includes material destination assignment via grade control polygons — this is how trucks know where to tip based on material classification. Available in lightweight OEM-agnostic systems as well as enterprise platforms.",
			},
			{
				"title":       "Advanced (enterprise FMS — Tier 1 operations)",
				"description": "Dynamic dispatch optimisation, algorithmic truck-to-loader assignment, match factor management, real-time blast movement integration with grade control (post-blast polygon correction before loading), AHS interface. Requires significant infrastructure investment and operational discipline to extract value.",
			},
			{
				"title":       "Vendor landscape by tier",
				"description": "Tier 1 large operations: MineStar Fleet and Modular Mining DISPATCH dominate — deep integration, high capability, purpose-built for complex multi-loader environments. Mid/lower tier: OEM-agnostic lightweight systems are winning — faster to implement, lower cost, no single-vendor equipment dependency. The right choice depends on operation size, complexity, and maturity.",
			},
			{
				"title":       "The risk of buying ahead of your capability",
				"description": "Deploying advanced dispatch on an operation that has not mastered core data quality is one of the most common — and costly — FMS failures. The algorithm produces optimal assignments based on the data it receives. Poor position data and unreliable delay codes produce confidently wrong dispatch decisions.",
			},
		},
	})

	// ── Section 2: System Requirements ─────────────────────────────────────
	s2 := createSection(database, module.ID, "What the System Must Know", 1)

	createSlide(database, s2.ID, "What the System Must Know", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "What the System Must Know",
		"subtitle": "Defining FMS information requirements before selecting technology",
		"bullets": []string{
			"The three fundamental data requirements of every FMS",
			"How perception-based sensing satisfies those requirements",
			"From the vehicle edge to the back office",
		},
	})

	createSlide(database, s2.ID, "The Three Fundamental Requirements", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Three Fundamental Requirements",
		"introText": "Before selecting sensors or platforms, define what the system must know. Every FMS capability — from basic shift reporting to dynamic dispatch — is built on three data requirements being satisfied continuously and automatically.",
		"items": []map[string]interface{}{
			{
				"title":       "1. Position and heading",
				"description": "Sub-metre accuracy for every mobile asset, continuously. GNSS (L1/L2, RTK for sub-metre accuracy) is the baseline. Map-matching against a high-definition mine model corrects for GNSS error in degraded zones (pit walls, high berms, underground). Without reliable position, everything downstream is compromised.",
			},
			{
				"title":       "2. Activity state",
				"description": "Is this machine loading, hauling loaded, queuing, hauling empty, dumping, on a delay, or parked? Activity state must be inferred automatically — not entered by the operator. State inference uses engine load, suspension weight, proximity to known locations, and motion vectors. Target: >95% automatic classification rate. The remaining <5% enters a manual classification queue; if not actioned within shift, records default to 'unknown' — corrupting delay analysis and the utilisation metrics built on it. This is the same failure mode as poor delay code discipline.",
			},
			{
				"title":       "3. Payload",
				"description": "How much material is on board? Onboard payload monitoring (suspension pressure sensors) provides real-time tonnes per truck per cycle. A well-commissioned system achieves ±2–3% accuracy — the target specification. ±5% is the acceptable upper bound triggering recalibration; operating above this threshold degrades production accounting and undermines the Measure step. High-value capability — enables payload compliance, tonne accounting, and tyre load accumulation tracking.",
			},
		},
	})

	createSlide(database, s2.ID, "Perception as the Enabling Data Layer", "PROCESS", 2, map[string]interface{}{
		"type":        "process",
		"heading":     "Perception as the Enabling Data Layer",
		"description": "Requirements define what the system must know. Perception technologies are selected to satisfy those requirements — not the other way around. Each requirement maps to one or more sensing modalities.",
		"steps": []map[string]interface{}{
			{
				"title":       "Position and heading",
				"description": "GNSS with RTK correction for sub-metre accuracy. IMU for dead-reckoning in GNSS-degraded areas. Map-matching against a continuously updated mine model. Multiple GNSS constellations (GPS, GLONASS, Galileo) for redundancy in pit environments with high sky obstruction.",
			},
			{
				"title":       "Activity state inference",
				"description": "Engine load and RPM, suspension pressure (loaded vs empty), proximity to known loader and dump positions, speed and motion vectors. A rules-based state machine or ML classifier processes these inputs to determine activity state without operator interaction. Accuracy target: >95% auto-classification rate.",
			},
			{
				"title":       "Payload measurement",
				"description": "Front and rear suspension strut pressure sensors measure hydraulic pressure, which is converted to tonnes via a calibrated load curve. Accuracy: ±2–5% of payload under normal conditions. Requires recalibration after tyre or strut changes. Some platforms use load pins in the body hinge as an alternative.",
			},
			{
				"title":       "Environmental context",
				"description": "Road network digitisation, loader and dump positions, grade profiles, and geofenced zones complete the picture. Static data updated as the mine advances. Combined with the real-time sensor streams, this gives the FMS the full operational context needed for production recording and dispatch decisions.",
			},
		},
	})

	createSlide(database, s2.ID, "Edge to Back Office — Data Pipeline", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Edge to Back Office — Data Pipeline",
		"introText": "FMS generates high volumes of time-series data from every vehicle on every shift. The network architecture that moves this data from machine to back office is a first-class design constraint — not an afterthought. Most FMS performance failures are communications failures.",
		"items": []map[string]interface{}{
			{
				"title":       "Latency requirements",
				"description": "Dispatch updates must reach trucks within 1–3 seconds to be operationally useful. Onboard safety enforcement (speed limits, exclusion zones) must operate independently of back office connectivity — latency-critical logic runs at the edge.",
			},
			{
				"title":       "Bandwidth requirements",
				"description": "Basic telemetry (GNSS, state, payload, engine data) requires 5–20 kbps per vehicle. Add video or high-frequency sensor data and this grows to 1–10 Mbps. Network design must account for peak fleet density zones — the pit ramp at shift change is the stress test.",
			},
			{
				"title":       "Network topology",
				"description": "802.11 Wi-Fi mesh (access points on pit infrastructure) is dominant on legacy and brownfield sites. New greenfield deployments increasingly use private LTE or 5G as primary, with Wi-Fi at fixed infrastructure points. Legacy sites retain 400 MHz UHF radio as a fallback layer for critical messages — new installations typically design this out in favour of LTE redundancy. Redundancy is non-negotiable regardless of approach.",
			},
			{
				"title":       "Edge vs back office processing",
				"description": "Safety enforcement, onboard display, and immediate state logging run at the edge — they cannot depend on a network connection. Analytics, optimisation, reporting, and integration with mine systems run in the back office. The architecture must be explicit about which functions live where.",
			},
		},
	})

	createSlide(database, s2.ID, "FMS-Equipped Haul Truck", "IMAGE", 4, map[string]interface{}{
		"type":        "image",
		"heading":     "FMS-Equipped Haul Truck",
		"description": "Every truck in a modern fleet carries a suite of sensors that feed the FMS. The GNSS antenna sits on the cab roof; suspension pressure sensors are on the hydraulic struts; the CAN bus runs from the engine ECU through the cab; the operator's touchscreen display faces the driver; and a network antenna maintains the data link back to the back office. All of this is invisible to an observer on the ground — but it is generating data on every cycle.",
		"imageUrl":    "https://assets.miningacademy.com/images/fms-equipped-haul-truck-sensors.jpg",
	})
	createSlide(database, s2.ID, "FMS On-Board Architecture", "DIAGRAM", 5, map[string]interface{}{
		"type":        "diagram",
		"heading":     "FMS On-Board Architecture — Sensor Stack to Back Office",
		"description": "Every FMS installation follows the same architecture: perception sensors on the vehicle generate raw data, onboard computing classifies and transmits it, and the back office aggregates, records, and dispatches. Understanding this stack determines where failures occur and what to fix first.",
		"diagramType": "architecture",
		"diagramData": map[string]interface{}{
			"sensors": []map[string]interface{}{
				{
					"name":        "GNSS / Positioning",
					"description": "L1/L2 dual-frequency with RTK correction. IMU for dead-reckoning in GNSS-degraded zones.",
					"functions":   []string{"Position", "Heading", "Speed", "Map-match"},
				},
				{
					"name":        "Payload Sensors",
					"description": "Hydraulic strut pressure sensors on suspension, converted to tonnes via calibrated load curve. Target: ±2–3%.",
					"functions":   []string{"Cycle payload", "Overload detection", "Tyre load tracking"},
				},
				{
					"name":        "Engine / OBD Interface",
					"description": "CAN bus data from engine ECU: RPM, load %, fault codes, fuel consumption rate.",
					"functions":   []string{"Engine hours", "Fault codes", "Fuel consumption"},
				},
				{
					"name":        "Operator Display",
					"description": "Cab touchscreen showing assignment, speed limit, delay codes, payload, and back-office messages.",
					"functions":   []string{"Delay coding", "Assignment ack", "Operator alerts"},
				},
				{
					"name":        "Network Infrastructure",
					"description": "Private LTE/5G (greenfield) or Wi-Fi mesh (brownfield). Latency target: <3s for dispatch updates.",
					"functions":   []string{"Data upload", "Dispatch delivery", "Time sync"},
				},
			},
			"coreAlgorithm": map[string]interface{}{
				"stages": []string{
					"Localise\n(GNSS + map-match)",
					"Classify\n(activity state)",
					"Record\n(cycle + events)",
					"Dispatch\n(assignment optimisation)",
				},
			},
		},
	})

	// ── Section 3: Production Records & Data (Measure) ──────────────────────
	s3 := createSection(database, module.ID, "Production Records and Data", 2)

	createSlide(database, s3.ID, "Production Records and Data", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Production Records and Data",
		"subtitle": "Measure — the foundation of operational accountability",
		"bullets": []string{
			"What gets recorded in every FMS cycle",
			"Shift production reports and the metrics that matter",
			"Delay analysis — turning event data into improvement actions",
		},
	})

	createSlide(database, s3.ID, "What Gets Recorded", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What Gets Recorded Every Cycle",
		"introText": "Every truck cycle generates a structured data record. The quality of downstream analytics — shift reports, delay analysis, mine planning feedback, maintenance triggers — is entirely dependent on the completeness and accuracy of this per-cycle record.",
		"items": []map[string]interface{}{
			{
				"title":       "Cycle definition",
				"description": "One complete truck cycle: spot at loader → load → haul loaded → spot at dump → dump → haul empty → return to loader. Each phase is timestamped. Phase durations are the raw inputs to all productivity analysis.",
			},
			{
				"title":       "Per-cycle data captured",
				"description": "Loader identity, dump destination, material classification (ore/waste/grade), payload (tonnes), haul distance, cycle time broken down by phase, operator identity, fuel consumed, and any delay events within the cycle.",
			},
			{
				"title":       "Drill cycle data",
				"description": "For drills: hole ID, pattern ID, design depth vs actual depth, penetration rate, rotary pressure, pulldown force, rotation speed, hole deviation, bit type, and re-drill events. Advanced systems use this data for MWD (Measurement While Drilling) — inferring lithology and rock strength from drilling parameters to improve blast design. Drill FMS data feeds blast design review, drill performance benchmarking, and ultimately fragmentation quality.",
			},
			{
				"title":       "Accuracy and completeness",
				"description": "A cycle record with a missing payload, wrong destination, or misclassified material classification is worse than no record — it actively misleads production accounting and ore tracking. Data quality auditing (completeness rate, payload capture rate, auto-classification rate) should be tracked as KPIs in their own right.",
			},
		},
	})

	createSlide(database, s3.ID, "Shift Production Reports", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Shift Production Reports",
		"introText": "The shift report is the primary accountability document for operational performance. A well-structured FMS shift report answers the following questions automatically — no manual compilation required.",
		"items": []map[string]interface{}{
			{
				"title":       "Tonnes moved vs target",
				"description": "Total BCM and wet metric tonnes by material classification and destination, variance to shift target as percentage. The headline — but it hides the story. Two shifts hitting the same tonnage target can look identical here while having completely different efficiency profiles.",
			},
			{
				"title":       "Cycles per loader per hour",
				"description": "Isolates loader productivity from fleet availability. A loader producing 18 cycles/hr with 85% utilisation tells a different story than one at 22 cycles/hr with 60% utilisation. This is the metric that identifies whether the loader or the fleet is the constraint.",
			},
			{
				"title":       "Truck utilisation breakdown",
				"description": "Productive time vs queuing time vs delay time vs maintenance time, per truck per shift. The decomposition of non-productive time is where the improvement opportunities live — not the headline utilisation number.",
			},
			{
				"title":       "Variance to design",
				"description": "Actual cycle time vs design cycle time by haul road segment. Systematic overruns on a specific road segment point to road condition, grade, or congestion issues. This feeds directly into the road maintenance and mine planning workflow.",
			},
		},
	})

	createSlide(database, s3.ID, "Delay Analysis and Classification", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Delay Analysis and Classification",
		"introText": "Delay data is the highest-value, lowest-quality dataset in most FMS deployments. The technology to capture it is straightforward. The discipline to keep it accurate and granular over time is the hard part.",
		"items": []map[string]interface{}{
			{
				"title":       "Delay taxonomy",
				"description": "Equipment breakdowns and repairs reduce Mechanical Availability — these are downtime events. Scheduled delays (shift change, crib, planned maintenance) and Standby events (blast hold, weather, operator shortage) reduce Utilisation of Available Time — the machine is available but not operating. Getting this distinction right is critical: it determines whether the improvement owner is maintenance (downtime) or operations (standby/schedule). The taxonomy here feeds the TUM hierarchy covered in the Fleet Visibility section.",
			},
			{
				"title":       "Automatic vs manual classification",
				"description": "FMS auto-classifies delays using state inference — a stationary truck with engine running in no queue position is likely on an unscheduled delay. Operator confirmation adds specificity. The hybrid approach (auto-detect + operator confirm within 2 minutes) achieves the best accuracy without relying on manual entry.",
			},
			{
				"title":       "What good delay codes look like",
				"description": "Bad: 'Mechanical delay'. Good: 'Engine fault — sustained high exhaust gas temperature, Grade 4 haul road, ambient >38°C'. The specificity required to drive action is also the specificity that prevents gaming. Code design should be done by operations and maintenance engineers, not IT.",
			},
			{
				"title":       "Rolling pattern analysis",
				"description": "Single-shift delay data is noise. Thirty-day rolling aggregates reveal the systemic patterns: the loader with disproportionate respot delays; the road segment that creates queue build-up at the same time each shift; the operator who consistently logs 'other' delays. These patterns are invisible in day-to-day reporting.",
			},
		},
	})

	// ── Section 4: Fleet Visibility & Management (Manage) ───────────────────
	s4 := createSection(database, module.ID, "Fleet Visibility and Management", 3)

	createSlide(database, s4.ID, "Fleet Visibility and Management", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Fleet Visibility and Management",
		"subtitle": "Manage — giving supervisors the tools to act on live data",
		"bullets": []string{
			"Real-time supervisor tools and decision support",
			"Payload and fuel management",
			"Utilisation framework — connecting FMS to TUM",
		},
	})

	createSlide(database, s4.ID, "Real-Time Supervisor Tools", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Real-Time Supervisor Tools",
		"introText": "The Manage layer is where data becomes decisions. A supervisor with a real-time FMS dashboard can intervene in situations that would otherwise compound silently through the shift. The technology enables this — but only if the supervisor knows how to read it.",
		"items": []map[string]interface{}{
			{
				"title":       "Live fleet map",
				"description": "Every machine positioned in real time on the mine map with activity state colour coding. A supervisor can see queue build-up at a loader, a truck sitting on an undeclared delay, or a haul road congestion point developing — and act before it cascades into a production loss.",
			},
			{
				"title":       "Delay management console",
				"description": "Active delays by machine, duration, and classification. Sorted by impact — the longest-running unresolved delay at the top. The supervisor's job is to close delays as fast as possible; the console makes the priority ordering visible. Average delay response time is a supervisor performance metric.",
			},
			{
				"title":       "Shift progress vs target",
				"description": "Running tonnes against shift target, projected end-of-shift total based on current rate, and variance trend. A supervisor who can see at 10:00am that they are tracking 8% below target has options. One who finds out at shift change has none.",
			},
			{
				"title":       "Shift change management",
				"description": "Shift handover is one of the largest single sources of production loss — typically 20–40 minutes of fleet dead time per shift change, twice per day. FMS enables pre-positioned truck assignments, supervisor handover notes linked to live machine states, and automatic recording of changeover duration.",
			},
		},
	})

	createSlide(database, s4.ID, "Payload and Fuel Management", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Payload and Fuel Management",
		"introText": "Payload and fuel are the two highest-leverage operational levers accessible through FMS data. Both are invisible without instrumentation and both reward consistent management over time.",
		"items": []map[string]interface{}{
			{
				"title":       "Payload consistency over average payload",
				"description": "The target is not the highest possible payload — it is the most consistent payload. Tyre load-life degrades non-linearly with overloading: approximately 30% reduction per 10% overload is a commonly cited approximation from OTR manufacturer load-life curves, though exact figures are tyre-construction and haul-profile specific. For a 59/80R63 (63-inch rim) OTR tyre at ~$80,000 on an ultra-class truck, the financial stakes of systematic overloading are significant. Variance is the enemy.",
			},
			{
				"title":       "The cost of underloading",
				"description": "A truck 10% under target delivers 10% less material per cycle. To compensate, the fleet needs 11% more cycles — consuming proportionally more fuel, tyres, maintenance hours, and driver time. Chronic underloading is invisible in headline tonnage reports until you cross-reference payload data.",
			},
			{
				"title":       "Fuel and idle time",
				"description": "A large haul truck consumes ~20–30 L/hr at idle versus 80–120 L/hr under load. FMS classifies idle time by category (queue, shift change, refuel, mechanical, operator-initiated) — the decomposition tells you where to intervene. Queue idle is a dispatch problem. Shift change idle is a schedule problem. Operator idle is a training problem.",
			},
			{
				"title":       "Grade and speed interaction",
				"description": "Fuel burn increases non-linearly with road grade and speed. FMS position data linked to surveyed grade profiles identifies specific road segments with disproportionate fuel cost — direct input to road maintenance prioritisation. A 1 percentage point reduction in ramp grade (e.g., 8% to 7%) on a loaded haul segment can reduce fuel consumption on that segment by 5–8% depending on truck class and load. At fleet scale this compounds significantly.",
			},
		},
	})

	createSlide(database, s4.ID, "Utilisation Framework — FMS and TUM", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Utilisation Framework — FMS and TUM",
		"introText": "The Time Utilization Model (TUM) is the standard mining industry framework for decomposing equipment time: Calendar Time → Scheduled Time → Available Time → Operating Time → Productive Time. FMS is the primary data source that makes this hierarchy measurable with precision.",
		"items": []map[string]interface{}{
			{
				"title":       "FMS as the TUM data source",
				"description": "Mechanical Availability (Available Time / (Available Time + Downtime)) and Utilisation of Available Time (Operating Time / Available Time) are both derived from FMS-recorded time. Equipment breakdowns reduce Mechanical Availability; unproductive operating time (queue, idle, extended handover, standby events like blast holds) reduces Utilisation. Without FMS, both require manual logging — lower accuracy, no cycle-level decomposition.",
			},
			{
				"title":       "Where FMS improvement attacks the TUM hierarchy",
				"description": "Utilisation of Available Time losses — queue time, extended shift change, inefficient routing — are all within FMS Manage and Optimise capability. A structured FMS programme directly improves the Utilisation component. Shift change duration (a Scheduled Time loss) is also partially addressable via FMS-supported handover management.",
			},
			{
				"title":       "Where FMS cannot help",
				"description": "Physical Availability losses — breakdowns and scheduled maintenance — are in the domain of CMMS and maintenance strategy. FMS records them with precision but cannot prevent them. FMS-CMMS integration bridges these two domains, using FMS operational data to trigger condition-based maintenance decisions.",
			},
			{
				"title":       "Operating efficiency — the dimension utilisation misses",
				"description": "The TUM hierarchy also captures operating efficiency: how close to design performance the machine runs when it is operating. A truck at full utilisation but 10% slower cycle times (road condition, loaded speed compliance, operator behaviour) is invisible in utilisation reporting — but visible when actual cycle time is compared against design. This is the metric that rewards excellence, not just presence.",
			},
		},
	})

	createSlide(database, s4.ID, "Key FMS Performance Metrics", "DIAGRAM", 4, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Key FMS Performance Metrics",
		"description": "These are the metrics FMS makes measurable. Understanding what each formula captures — and who owns it — is fundamental to using FMS data correctly.",
		"diagramType": "kpi-map",
		"diagramData": map[string]interface{}{
			"ratios": []map[string]interface{}{
				{
					"name":        "Physical Availability (PA)",
					"formula":     "Available Time / (Available Time + Downtime)",
					"owner":       "Maintenance",
					"description": "Percentage of scheduled time the machine is mechanically ready to operate. Losses here are breakdowns and planned maintenance — outside FMS's direct control but precisely measured by it.",
				},
				{
					"name":        "Utilisation of Available Time (UAT)",
					"formula":     "Operating Time / Available Time",
					"owner":       "Operations",
					"description": "Of the time the machine is available, how much is actually working? Queue, extended shift change, standby (blast hold, weather), and operator idle all reduce UAT — these are the FMS Manage-layer targets.",
				},
				{
					"name":        "Match Factor (MF)",
					"formula":     "(N trucks × loader cycle time) / truck cycle time",
					"owner":       "Operations / Dispatch",
					"description": "Ratio of truck supply to loader demand. Target: MF slightly below 1.0 (slight over-trucking keeps the loader fed). MF < 0.85: loader starving. MF > 1.15: trucks queuing.",
				},
				{
					"name":        "Tonnes per Effective Equipment Hour (TPEEH)",
					"formula":     "Total Tonnes / Effective Equipment Hours",
					"owner":       "Planning / Management",
					"description": "Comprehensive productivity metric combining availability, utilisation, and payload efficiency. Benchmarked against mine plan design parameters. The headline metric for FMS programme performance review.",
				},
			},
		},
	})

	createSlide(database, s4.ID, "Day in the Life: Shift Supervisor", "PROCESS", 5, map[string]interface{}{
		"type":        "process",
		"heading":     "Day in the Life: Shift Supervisor",
		"description": "A production supervisor's 12-hour shift using FMS as their primary management tool. The FMS dashboard is open from first minute to last. Every decision is data-referenced.",
		"steps": []map[string]interface{}{
			{
				"title":       "06:00 — Pre-shift handover",
				"duration":    "30 min",
				"description": "Review the outgoing supervisor's FMS handover report: night shift tonnes vs plan, equipment availability list, active delays, and any machines on maintenance hold. The FMS summary takes 5 minutes; the conversation takes the remaining 25.",
			},
			{
				"title":       "06:30 — Crew allocation",
				"duration":    "30 min",
				"description": "Confirm operator-to-machine assignments against the FMS-reported availability list. Any machines still logging a prior delay code get investigated before the operator boards. A machine 'available' in FMS that the outgoing crew flagged as rough needs a decision.",
			},
			{
				"title":       "07:00 — Pit walkthrough",
				"duration":    "60 min",
				"description": "Physical inspection of loader positions, road conditions, and blast prep areas — cross-referenced with what FMS is showing. Discrepancies between the physical state and FMS state (wrong loader position, missing truck, unlogged delay) are corrected now, before they corrupt the shift report.",
			},
			{
				"title":       "08:00 — First production check",
				"duration":    "15 min",
				"description": "Review first-hour actual production rate vs target. Check payload compliance distribution — any trucks consistently over or under target get flagged. A supervisor who spots a systematic payload issue at 08:00 has 9 hours to fix it. One who sees it at shift change has nothing.",
			},
			{
				"title":       "09:00 — Delay audit",
				"duration":    "30 min",
				"description": "Review all delays logged in the first 2 hours. Any coded as 'other' or 'operator hold' get investigated — are these correctly classified? Delay code discipline is a supervisor accountability metric. Fixing it now keeps the shift data clean.",
				"warnings": []string{
					"Delay code quality degrades within 6 months if supervisors don't enforce it weekly",
				},
			},
			{
				"title":       "12:00 — Mid-shift production review",
				"duration":    "20 min",
				"description": "Formal mid-shift check: shift-to-date tonnes against plan, utilisation stats by machine, match factor trend by loader. If tracking >5% behind target, escalate with a root cause — not just the number. The dispatcher and planner both get the same FMS data; the supervisor adds context.",
			},
			{
				"title":       "13:00 — Blast coordination",
				"duration":    "Variable",
				"description": "Pre-blast standby management: all trucks in the blast radius correctly logged as blast hold delay (Standby — not a breakdown). Post-blast: confirm grade control has updated polygons before loading resumes. The polygon update workflow is the supervisor's gate — it doesn't happen automatically.",
			},
			{
				"title":       "17:00 — End-of-shift handover",
				"duration":    "30 min",
				"description": "Formal handover using the FMS shift summary as the primary document. Equipment states, active delays, maintenance flags, tonnes vs plan, outstanding data quality issues. The incoming supervisor inherits your data quality — leave it accurate.",
			},
		},
	})

	// ── Section 5: Dispatch & Optimisation (Optimise — Advanced) ────────────
	s5 := createSection(database, module.ID, "Dispatch and Optimisation", 4)

	createSlide(database, s5.ID, "Dispatch and Optimisation", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Dispatch and Optimisation",
		"subtitle": "Optimise — the advanced capability layer of FMS",
		"bullets": []string{
			"Dynamic assignment as the top of the FMS value ladder",
			"How dispatch algorithms work and what they optimise for",
			"Match factor, queue management, and key KPIs",
		},
	})

	createSlide(database, s5.ID, "The Dispatch Problem", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Dispatch Problem",
		"introText": "Dynamic dispatch is the algorithmic engine at the top of the FMS value ladder. It solves a continuous assignment problem: given N trucks at various positions with various payloads, assign each to one of M loaders to maximise total tonnes moved — subject to safety, road, and destination constraints.",
		"items": []map[string]interface{}{
			{
				"title":       "Why it requires the foundations first",
				"description": "The dispatch engine optimises against the data it receives. Poor GNSS positions produce wrong haul time estimates. Missing payloads break the match factor calculation. Unreliable delay data creates phantom availability. A dispatch algorithm on poor foundations produces confidently wrong assignments.",
			},
			{
				"title":       "Continuous re-optimisation",
				"description": "The optimal assignment found at t=0 is typically suboptimal by t=60 seconds as trucks move, queues change, and equipment fails. Modern FMS dispatch engines re-solve the assignment problem every 30–120 seconds — the interval must be shorter than the time for a truck to pass its last practical assignment decision point (the last intersection before committing to a haul segment). Beyond that point, a new assignment cannot be acted on until the next cycle.",
			},
			{
				"title":       "Objective function",
				"description": "Maximise total tonnes to destination while meeting grade blend requirements, minimising total truck-hours (including queuing), and respecting road capacity constraints. Advanced systems add tyre life accumulation and fuel burn to the objective — true cost-per-tonne optimisation.",
			},
			{
				"title":       "Constraints",
				"description": "Road segment capacity (maximum simultaneous trucks), destination availability (crusher blocked, dump face full), equipment class compatibility (truck size vs loader type), and safety exclusions (blast zones, maintenance areas, AHS corridors).",
			},
		},
	})

	createSlide(database, s5.ID, "How Dispatch Algorithms Work", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "How Dispatch Algorithms Work",
		"introText": "Most commercial FMS dispatch engines combine linear programming (LP) for global optimisation with heuristic rules that encode operational knowledge the LP cannot capture.",
		"items": []map[string]interface{}{
			{
				"title":       "Linear programming formulation",
				"description": "Decision variables: binary assignment of each available truck to each loader. Objective: minimise expected queue time plus haul time across all assignments. Constraints: one assignment per truck, loader capacity limits, road constraints. Solved using simplex or interior-point methods at each re-solve interval.",
			},
			{
				"title":       "Heuristic overrides",
				"description": "Rule-based dispatch for specific conditions: priority loaders for ore blend compliance, blast hold management, crusher queue limits, emergency re-assignments. The heuristics encode operational knowledge — the LP handles scale and combinatorial complexity.",
			},
			{
				"title":       "Prediction quality drives assignment quality",
				"description": "The LP optimises against predicted haul times. Those predictions are only as good as the historical cycle data feeding them. A fleet with consistent cycle time data gets better assignments than one with noisy records — another compounding return on data quality investment.",
			},
			{
				"title":       "ML extensions",
				"description": "Emerging capability: ML models trained on historical cycle data predict haul times more accurately than static distance-based estimates, accounting for road condition, payload, ambient temperature, and operator behaviour patterns. Better predictions produce better assignments — but the training data requirement is significant.",
			},
		},
	})

	createSlide(database, s5.ID, "Match Factor and Queue Management", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Match Factor and Queue Management",
		"introText": "Match factor is the primary design parameter for dispatch optimisation. It governs the balance between loader idle time and truck queue time — both of which destroy productivity, but in different ways and with different cost profiles.",
		"items": []map[string]interface{}{
			{
				"title":       "Match factor definition and preferred operating point",
				"description": "Match Factor = (N trucks × loader cycle time) / truck cycle time. MF = 1.0 means truck arrivals perfectly fill loader capacity. The preferred operating point is MF slightly below 1.0 — marginal over-trucking, meaning there are always more trucks available than the loader strictly needs. This ensures the loader is never waiting. The asymmetry: loader idle costs $3,000–8,000/hr (rope shovel class); each queuing truck costs $300–600/hr. 'Bias toward the loader' means bias toward loader uptime — achieved by slight over-trucking, not by running MF above 1.0.",
			},
			{
				"title":       "Under-trucked (MF < 0.85) — the higher cost failure",
				"description": "The loader idles between trucks. Large rope shovel and electric mining shovel class equipment costs $3,000–8,000/hr to operate; hydraulic excavators typically $1,500–3,000/hr. Loader starvation is the highest-cost dispatch failure mode. A dispatch system should be biased toward slight over-trucking — queue time is recoverable; loader idle time is not. Thresholds are indicative and fleet-size dependent.",
			},
			{
				"title":       "Over-trucked (MF > 1.15)",
				"description": "Trucks queue at the loader. Queue time is dead time — consuming wages, fuel, and tyre life while producing nothing. The dispatch system re-assigns queuing trucks to alternative loaders or secondary tasks (fuel bay, maintenance access). Queue length is visible in real time on the supervisor dashboard.",
			},
			{
				"title":       "Dynamic match factor",
				"description": "Match factor shifts constantly as trucks break down, haul distances change (pit deepening), and loader availability shifts. Real-time dispatch continuously adjusts. Match factor by loader per hour, trended over a shift, is the primary diagnostic tool for dispatch performance review.",
			},
		},
	})

	createSlide(database, s5.ID, "Key Dispatch KPIs", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Key Dispatch KPIs",
		"introText": "These KPIs are the language of dispatch performance. The numbers matter — but understanding what drives each metric is what separates engineers who can improve performance from those who can only report it.",
		"items": []map[string]interface{}{
			{
				"title":       "Truck cycle time (decomposed)",
				"description": "Total cycle time broken into phases: spot + load + haul loaded + spot at dump + dump + haul empty + queue. A 60-minute cycle with 12 minutes of queue means 20% of capacity is consumed by queue. The decomposition points to the intervention — not the total.",
			},
			{
				"title":       "Loader utilisation (%)",
				"description": "Percentage of scheduled operating time the loader is loading a truck. Target >85%. Losses from truck gaps, shift change, respot, and blast hold. This is the metric the dispatch system is primarily optimising — maximising loader productive time.",
			},
			{
				"title":       "Tonnes per effective equipment hour (TPEEH)",
				"description": "Production rate adjusted for all delay losses. Combines truck availability, utilisation, and payload efficiency into a single number. The most comprehensive productivity metric. Benchmarked against design parameters established during mine planning.",
			},
			{
				"title":       "Average payload vs target with variance",
				"description": "Mean payload as % of target and standard deviation. Target: 98–102% mean, standard deviation <5%. High mean with high variance indicates overloading events within otherwise good average performance — a tyre risk that headline average does not reveal.",
			},
		},
	})

	createSlide(database, s5.ID, "Day in the Life: FMS Dispatcher", "PROCESS", 5, map[string]interface{}{
		"type":        "process",
		"heading":     "Day in the Life: FMS Dispatcher",
		"description": "A dispatcher in a large open-cut mine operating a full enterprise FMS with automated dispatch. The dispatcher's job is not to manually assign trucks — it is to monitor the algorithm, intervene when it can't, and maintain the data quality that makes it work.",
		"steps": []map[string]interface{}{
			{
				"title":       "06:00 — System handover",
				"duration":    "20 min",
				"description": "Review overnight dispatch performance: average match factor by loader, total overrides made by the night dispatcher, any trucks flagged for maintenance, payload anomalies. Check that all machines logging on for day shift have valid positions — a truck logging on without a GNSS fix will receive wrong assignments.",
			},
			{
				"title":       "06:30 — Production ramp-up monitoring",
				"duration":    "60 min",
				"description": "Watch the fleet come online. Each loader reaching operating position and each truck completing its first dispatch assignment. The first 30 minutes reveal data quality problems: trucks not acknowledging assignments, loaders at wrong map positions, payload sensors reading zero. Fix these before they compound.",
			},
			{
				"title":       "07:30 — Match factor monitoring",
				"duration":    "Continuous",
				"description": "Primary dashboard: match factor live by loader. Any loader dropping below MF 0.85 for more than two consecutive cycles triggers an investigation — is the haul segment longer than designed? A truck broken down? Uncommonly slow loader cycle? The algorithm will re-assign, but the dispatcher needs to understand why.",
			},
			{
				"title":       "09:00 — Cycle time validation",
				"duration":    "30 min",
				"description": "Compare actual cycle times by haul segment against design. Any truck running >15% longer than design on a specific segment flags a road condition, grade, or operator behaviour issue — escalated to the supervisor, not handled by the dispatcher. The dispatcher identifies the pattern; the supervisor owns the fix.",
			},
			{
				"title":       "10:00 — Override review",
				"duration":    "15 min",
				"description": "Review override log from the past 3 hours. Any override made by a supervisor should have a reason. Pattern of overrides on a specific loader or haul segment may indicate an algorithm parameter issue — document it for the post-shift review with the FMS system coordinator.",
				"warnings": []string{
					"High override rates are a training signal, not normal dispatcher behaviour",
				},
			},
			{
				"title":       "12:00 — Blast hold management",
				"duration":    "Variable",
				"description": "Blast hold period: manually move all trucks in the exclusion zone to blast hold delay code. The algorithm pauses dispatch for those trucks; the dispatcher confirms each machine is correctly parked and logged. Post-blast, restore machines to available and monitor the algorithm's re-assignment of the fleet.",
			},
			{
				"title":       "14:00 — Productivity report",
				"duration":    "15 min",
				"description": "Run the shift-to-date productivity summary: tonnes vs plan, loader utilisation %, average match factor, payload compliance %. If trailing >5% against plan with 3 hours to run, escalate to the shift supervisor with the specific constraint — not just the gap.",
			},
			{
				"title":       "17:30 — Dispatcher handover",
				"duration":    "20 min",
				"description": "Detailed handover to the afternoon dispatcher: active dispatch state, any trucks approaching end-of-shift maintenance threshold, recurring issues observed, override patterns, and any FMS data quality issues unresolved. The afternoon dispatcher inherits the state — leave it clean.",
			},
		},
	})

	// ── Section 6: Integration with Mine Systems ────────────────────────────
	s6 := createSection(database, module.ID, "Integration with Mine Systems", 5)

	createSlide(database, s6.ID, "Integration with Mine Systems", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Integration with Mine Systems",
		"subtitle": "FMS as a hub — connected to grade control, maintenance, and autonomous haulage",
		"bullets": []string{
			"Grade control integration — FMS as a value chain tool",
			"CMMS integration — using FMS data to trigger maintenance",
			"AHS interface — managing mixed autonomous and manned fleets",
		},
	})

	createSlide(database, s6.ID, "Grade Control and Ore Tracking", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Grade Control and Ore Tracking",
		"introText": "Integrating FMS with grade control transforms it from a logistics tool into a mine value chain tool. Every truck load becomes a tracked material sample with a known source location, classification, and destination.",
		"items": []map[string]interface{}{
			{
				"title":       "Polygon-based material classification",
				"description": "Grade control geologists define ore/waste polygons on the mine surface model. FMS tracks each loader's position against these polygons to classify material being loaded in real time. Destination rules are auto-generated from material type — ore to ROM, waste to dump.",
			},
			{
				"title":       "Dispatch-level material tracking",
				"description": "Every tonne that leaves a loader is tracked to its dump destination by FMS — loader identity, material classification, tonnes, and destination recorded per cycle. This produces a real-time material movement log: what went where, and when. This is distinct from ore reconciliation (which compares dispatched material to assay results and is a post-event geological process) — FMS provides the movement data; the grade control system performs the reconciliation against it.",
			},
			{
				"title":       "Blast movement — the critical integration point",
				"description": "Post-blast, material polygons shift by 5–30m due to blast movement. If blast movement vectors are not applied to the grade control model before loading begins, trucks haul material classified against the pre-blast position — the highest-cost ore misclassification mechanism in open-cut mining. FMS-grade control integration must include blast movement update workflow.",
			},
			{
				"title":       "Dynamic destination updates",
				"description": "As grade control updates material boundaries during the shift (new blast data, in-pit sampling results), FMS updates routing rules dynamically. A truck mid-shift can be re-routed from ROM to waste dump based on revised grade interpretation — without supervisor intervention.",
			},
		},
	})

	createSlide(database, s6.ID, "CMMS Integration", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "CMMS Integration — FMS as a Maintenance Trigger",
		"introText": "FMS records engine hours, payload cycles, fault codes, and operating conditions for every machine on every shift. This is precisely the data a predictive maintenance system needs. Most operations have it but do not connect it.",
		"items": []map[string]interface{}{
			{
				"title":       "Automatic meter reading",
				"description": "FMS transmits real-time engine hours to CMMS, eliminating manual meter reading. Maintenance work orders trigger automatically at the correct threshold — not at the next shift change. Time between fault and work order creation drops from hours to minutes.",
			},
			{
				"title":       "Fault code-triggered maintenance",
				"description": "Engine and drivetrain fault codes from the machine's OBD system feed through FMS to CMMS in real time. A sustained high exhaust temperature fault auto-generates a CMMS inspection task and flags the truck for maintenance hold at the next available opportunity — before a hard failure.",
			},
			{
				"title":       "Condition-based suspension monitoring",
				"description": "Strut pressure readings on an empty truck should be consistent and within a known range. Anomalous empty-truck pressure readings indicate seal failure or nitrogen loss — detectable weeks before a hard strut failure if the data is being monitored. FMS-CMMS integration enables this early warning automatically.",
			},
			{
				"title":       "Tyre life integration",
				"description": "Payload per cycle multiplied by haul distance gives a tyre load accumulation metric far more predictive of tyre life than calendar time or hours alone. FMS-CMMS integration enables tyre life tracking on this basis — improving tyre change scheduling accuracy by 20–40% and reducing unexpected tyre failures.",
			},
		},
	})

	createSlide(database, s6.ID, "AHS Interface and Mixed Fleet Management", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "AHS Interface and Mixed Fleet Management",
		"introText": "Autonomous Haulage Systems change the FMS architecture fundamentally. The dispatch engine now manages vehicles with no human operator — with different safety envelopes, communication requirements, and performance characteristics.",
		"items": []map[string]interface{}{
			{
				"title":       "How FMS integrates with AHS",
				"description": "The FMS dispatch engine sends assignments to AHS trucks via a vehicle management interface (VMI) — the same assignment structure as manned dispatch, but the truck navigates autonomously using onboard path planning and obstacle avoidance. FMS manages fleet logistics; AHS manages vehicle control.",
			},
			{
				"title":       "Mixed fleet complexity",
				"description": "The most operationally demanding configuration is partial automation — some trucks AHS, some manned. Interaction rules must be encoded in the FMS: manned vehicles yield to AHS-controlled paths in designated zones; AHS trucks stop and request manual override when encountering unexpected manned behaviour.",
			},
			{
				"title":       "Productivity profile of AHS",
				"description": "AHS trucks typically run at 85–95% of manned truck individual cycle performance but achieve higher fleet productivity by eliminating shift change dead time, fatigue-related slowdowns, and technique variability. Compared against design targets or industry peer benchmarks (greenfield operations have no pre-AHS manned baseline), sites like Rio Tinto's Pilbara operations and Fortescue have reported 15–20% improvement in fleet tonnes per effective equipment hour. Brownfield and mixed-fleet results vary significantly — do not apply greenfield benchmarks to retrofit deployments.",
			},
			{
				"title":       "Infrastructure and discipline requirements",
				"description": "AHS requires network uptime >99.9%, high-definition mine maps updated daily, precisely surveyed road infrastructure, and strict exclusion zone management. A manned vehicle straying into an AHS exclusion zone triggers a fleet-wide stop in most systems. The operational discipline required is significantly higher than manned-only operations.",
			},
		},
	})

	// ── Section 7: Implementation & Change Management ────────────────────────
	s7 := createSection(database, module.ID, "Implementation and Change Management", 6)

	createSlide(database, s7.ID, "Implementation and Change Management", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Implementation and Change Management",
		"subtitle": "Why FMS programmes succeed and fail — and how to extract the value",
		"bullets": []string{
			"Implementation phases and sequencing",
			"The most common failure points",
			"Change management for the pit crew and supervisors",
		},
	})

	createSlide(database, s7.ID, "Implementation Phases", "PROCESS", 1, map[string]interface{}{
		"type":        "process",
		"heading":     "FMS Implementation Phases",
		"description": "A full FMS implementation from hardware to sustained performance takes 12–24 months. Rushing the foundation phases is the single most common cause of underperforming systems. Each phase builds the capability the next phase depends on.",
		"steps": []map[string]interface{}{
			{
				"title":       "Phase 1 — Infrastructure and Hardware → Value Ladder: Track (months 1–4)",
				"description": "Network infrastructure (Wi-Fi mesh, access points, repeaters), onboard hardware installation (GNSS antennas, computers, displays, payload sensors), back office server commissioning. This phase delivers Step 1 of the value ladder — Track. Data quality is determined here. Poor antenna placement and weak network coverage corrupt all downstream analytics. Expected productivity gain at phase completion: 5–8%.",
			},
			{
				"title":       "Phase 2 — Configuration and Baseline → Value Ladder: Measure (months 3–6)",
				"description": "Mine model import, road network digitisation, loader and dump point configuration, material classification rules, delay code taxonomy design. This phase delivers Step 2 — Measure. Operations engineering input is critical — the system is only as smart as the rules encoded into it. This phase should not be delegated to the vendor. Cumulative expected gain: 10–16%.",
			},
			{
				"title":       "Phase 3 — Advisory Mode and Data Validation → Value Ladder: Manage (months 5–8)",
				"description": "Supervisors use live FMS data to manage the shift — the system runs in advisory mode for dispatch. This delivers Step 3 — Manage. Validates data quality, builds operator familiarity, and establishes the baseline KPIs against which Phase 4 will be measured. Cumulative expected gain: 18–26%.",
			},
			{
				"title":       "Phase 4 — Automated Dispatch Go-Live → Value Ladder: Optimise (months 7–12)",
				"description": "Full automated dispatch activated — Step 4 of the value ladder. Requires intensive supervisor engagement in weeks 1–8; supervisors must understand algorithm logic to override correctly. Productivity typically dips 5–10% in the first 2 weeks before recovering and exceeding baseline. Full value ladder now active. The 30% headline is a realistic best-case for well-executed implementations; actual gains vary by site and are bounded by how thoroughly the first three steps were embedded.",
			},
			{
				"title":       "Phase 5 — Sustained Optimisation and Integration (months 12–24)",
				"description": "Fine-tune dispatch parameters against actual performance. Activate integrations (grade control blast movement, CMMS, survey). Begin feeding FMS actuals into mine planning. This phase generates the compounding long-term value — each percentage point improvement recurs every shift, indefinitely. Most operations under-invest in it.",
			},
		},
	})

	createSlide(database, s7.ID, "Common Failure Points", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Common Failure Points",
		"introText": "FMS implementation failures rarely occur in the technology. They occur in the operational discipline and organisational commitment required to extract value from it. The same patterns appear across operations of all sizes.",
		"items": []map[string]interface{}{
			{
				"title":       "Delay code discipline breakdown",
				"description": "Operators classify all unplanned delays as 'other'. Within 6 months, delay analysis is useless and the data foundation for Steps 2 and 3 collapses. Prevention: simple, meaningful taxonomy designed with operators; supervisor enforcement; weekly data quality audits visible to senior management.",
			},
			{
				"title":       "Supervisor override culture",
				"description": "Supervisors chronically override automated dispatch based on intuition, reverting to manual allocation. Often correct in the short term but wrong systemically — the algorithm sees the whole fleet; the supervisor sees their section. Prevention: train supervisors on dispatch logic; make override rates a visible management metric.",
			},
			{
				"title":       "Mine model staleness",
				"description": "The FMS mine model is not updated as the pit advances. Consequence: the dispatch engine routes trucks to mined-out faces, assigns loaders to polygons that no longer match actual dig positions, and produces assignments based on haul distances that no longer exist. Step 4 (Optimise) value is entirely lost — the algorithm is optimising against a map of a mine that no longer exists. All gains revert to manual dispatch. Prevention: defined update cadence (daily for active blast areas), clear ownership, and a production hold until the model is validated post-blast.",
			},
			{
				"title":       "Integration never activated",
				"description": "The FMS runs as an island — grade control, CMMS, and survey integrations remain disconnected for 12–18 months due to IT governance delays or undefined ownership. The system captures half its potential value. Prevention: integration architecture agreed and funded before hardware procurement, not after.",
			},
		},
	})

	// ── Section 8: Knowledge Check ──────────────────────────────────────────
	s8 := createSection(database, module.ID, "Knowledge Check", 7)

	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Mining Fleet Management Systems Assessment",
		PassMark:         80,
		MaxAttempts:      3,
		ShuffleQuestions: true,
		ShuffleOptions:   true,
		ShowResults:      true,
		SectionID:        &s8.ID,
		ModuleID:         &module.ID,
	}
	database.Create(&quiz)

	type option struct {
		text      string
		isCorrect bool
	}
	type question struct {
		question    string
		explanation string
		options     []option
	}

	questions := []question{
		{
			question:    "What are the three fundamental functions of every FMS, regardless of vendor or scale?",
			explanation: "Track (know where equipment is and what it's doing), Measure (turn tracking data into accurate production records and performance metrics), and Manage (give supervisors tools to act on that data). Dynamic dispatch is an advanced capability built on these foundations — not a core function of all FMS.",
			options: []option{
				{"Dispatch trucks, record production, send alerts", false},
				{"Track equipment, measure operational performance, manage site operations", true},
				{"Optimise haul routes, monitor payload, manage maintenance", false},
				{"Control autonomous vehicles, report utilisation, manage fuel", false},
			},
		},
		{
			question:    "A new operation is considering FMS. They have no existing GPS tracking and rely on radio dispatch. Which FMS value ladder step delivers the most immediate, lowest-risk productivity gain?",
			explanation: "Step 1 — Track — delivers 5–8% productivity improvement simply from GPS tracking and cycle counting, without changing any process. The gains come from eliminating lost trucks, untracked idle time, and inaccurate shift records. Higher steps require this foundation first.",
			options: []option{
				{"Step 4 — Optimise: deploy dynamic dispatch immediately for maximum gain", false},
				{"Step 1 — Track: GPS tracking and cycle counting delivers immediate gains on the lowest-risk foundation", true},
				{"Step 3 — Manage: real-time supervisor tools deliver the fastest ROI", false},
				{"Step 2 — Measure: accurate production records are the most urgent priority", false},
			},
		},
		{
			question:    "A blast engineer is reporting poor fragmentation in one zone despite the drill achieving design depth and penetration rate targets. A drill FMS MWD dataset is available. Which parameter is most likely to explain the fragmentation failure?",
			explanation: "Penetration rate and depth are basic drill progress metrics. The MWD parameters that reveal rock strength variation are pulldown force and rotary pressure — higher values in a zone of the same penetration rate indicate harder, more competent rock that requires a different explosive energy factor. If rotary pressure and pulldown force were not elevated despite poor fragmentation, the cause is more likely in the blast design (timing, burden, spacing) than the rock itself. MWD-derived lithology inference is the key FMS contribution to blast design quality.",
			options: []option{
				{"Penetration rate — faster drilling indicates softer rock that produces poor fragmentation", false},
				{"Hole deviation — misaligned holes reduce effective burden, producing coarse fragmentation", false},
				{"Rotary pressure and pulldown force — these reveal rock strength variation that penetration rate alone does not", true},
				{"Bit type — the wrong bit causes poor penetration and inadequate fragmentation energy transfer", false},
			},
		},
		{
			question:    "A match factor of 0.72 at a loader indicates what condition, and why is it the higher-cost failure mode?",
			explanation: "MF 0.72 is under-trucked — the loader is starving between truck arrivals. Rope shovel and electric mining shovel class equipment costs $3,000–8,000/hr to operate; hydraulic excavators $1,500–3,000/hr. Either way, loader idle is the highest-cost dispatch failure mode. A well-designed dispatch system is biased toward slight over-trucking (MF slightly below 1.0) — keeping the loader fed even at the cost of short truck queues.",
			options: []option{
				{"Over-trucked — trucks are building a queue, which is the highest-cost failure mode", false},
				{"Under-trucked — the loader is starving, which is the highest-cost failure mode due to loader operating costs of $3,000–8,000/hr (rope shovel class)", true},
				{"Perfectly balanced — MF of 0.72 is within the target range", false},
				{"Road congestion — trucks cannot reach the loader fast enough", false},
			},
		},
		{
			question:    "Why does payload variance matter more than low average payload?",
			explanation: "Tyre load-life degrades non-linearly with overloading — manufacturer empirical curves show the relationship is strongly superlinear, meaning even modest overloads compound disproportionate damage over time. The commonly cited approximation is ~30% tyre life reduction per 10% overload. High variance means frequent overloading events, each causing damage that the average payload figure never reveals. A consistently low average is suboptimal but predictable; high variance combines underloading losses with overloading damage.",
			options: []option{
				{"Average payload determines fuel burn per tonne, making it the primary cost driver", false},
				{"Variance causes inconsistent dispatch cycle times that confuse the optimisation algorithm", false},
				{"Overloading events within high variance cause non-linear tyre and structural degradation that average payload does not reveal", true},
				{"FMS systems can only accurately measure variance, not average payload", false},
			},
		},
		{
			question:    "What is the most critical integration point between FMS and grade control, and what goes wrong when it fails?",
			explanation: "The blast movement update — applying post-blast movement vectors to grade control polygons before loading begins. When this fails, trucks haul material classified against the pre-blast polygon position, which has shifted 5–30m. The result is ore hauled to waste or waste hauled to the crusher — the highest-cost ore misclassification mechanism in open-cut mining.",
			options: []option{
				{"Real-time mill feed rate reporting", false},
				{"Blast movement vectors applied to grade control polygons before loading begins", true},
				{"CMMS work order generation for drill equipment", false},
				{"Truck speed compliance in grade control exclusion zones", false},
			},
		},
		{
			question:    "An operation has deployed a full enterprise FMS with dynamic dispatch. Delay analysis data quality has degraded over 6 months — operators are logging 80% of delays as 'other'. What is the most likely operational consequence?",
			explanation: "The dispatch algorithm optimises against predicted cycle times derived from historical data. Delay code breakdown means the system cannot accurately predict when machines will be available — degrading assignment quality. The compounding effect is that the operation loses the value of Steps 2 and 3 (Measure, Manage), undermining the foundation that Step 4 (Optimise) depends on.",
			options: []option{
				{"The FMS network will become congested with unclassified delay messages", false},
				{"Dispatch assignment quality degrades as the algorithm loses accurate availability predictions, undermining the foundation that optimisation depends on", true},
				{"The CMMS integration will stop receiving fault codes", false},
				{"Payload monitoring accuracy will decrease", false},
			},
		},
		{
			question:    "A supervisor is consistently overriding automated dispatch assignments in their section. What is the most effective intervention?",
			explanation: "Override rates should be tracked and visible as a management metric. High override rates are typically a training issue — the supervisor does not understand why the algorithm made its decision, so they revert to intuition. Training on dispatch logic and making override rates visible to management reduces this. Removing override access damages operational trust and removes a necessary safety valve.",
			options: []option{
				{"Remove override access from supervisor accounts to enforce algorithm compliance", false},
				{"Track override rates, train supervisors on dispatch algorithm logic, and make override frequency a visible management metric", true},
				{"Adjust dispatch parameters to match the supervisor's historical allocation preferences", false},
				{"Override behaviour is expected and normal during early FMS deployment — no action needed", false},
			},
		},
		{
			question:    "A 60-truck fleet has Physical Availability of 91%, Utilisation of Available Time of 89%, and trucks are averaging 108-minute cycles against a 100-minute design. Shift tonnage is 4% below target. What is the primary constraint and what FMS action addresses it?",
			explanation: "Physical Availability (91%) and Utilisation (89%) are both healthy — the fleet is available and working. The constraint is operating efficiency: trucks are running 8% slower than design cycle time, which compounds across 60 trucks to produce the tonnage shortfall. The FMS action is to compare actual cycle times against design by road segment to identify where the slowdown is occurring — road condition, grade, or operator behaviour on a specific haul segment.",
			options: []option{
				{"Physical Availability — schedule additional maintenance to reduce breakdown frequency", false},
				{"Utilisation — reduce queue time at loaders by adjusting dispatch assignments", false},
				{"Operating efficiency — trucks are running 8% below design cycle time; investigate cycle time vs design by road segment", true},
				{"Payload — trucks are underloaded, requiring more cycles to hit target tonnage", false},
			},
		},
		{
			question:    "A 45-truck mixed-OEM open pit operation is deploying FMS for the first time. They have no existing GPS tracking, a limited IT team, and equipment from four different manufacturers. Which FMS approach is appropriate and why?",
			explanation: "This operation is at the Track step of the value ladder — their priority is visibility and basic production records, not dispatch optimisation. An OEM-agnostic lightweight system fits: lower implementation cost, no single-vendor equipment dependency, faster to deploy on a mixed fleet, and appropriately scoped to the operation's current capability maturity. Enterprise dispatch platforms (MineStar, Modular) are purpose-built for large, mature, complex multi-loader environments.",
			options: []option{
				{"Enterprise platform (MineStar Fleet or Modular DISPATCH) — only enterprise systems support mixed-OEM fleets", false},
				{"OEM-agnostic lightweight FMS — lower cost, faster implementation, no equipment vendor dependency, appropriate to current capability maturity", true},
				{"Single-OEM FMS from the dominant equipment manufacturer — standardise the fleet first", false},
				{"No FMS needed at 45 trucks — manual dispatch is sufficient at this scale", false},
			},
		},
	}

	for i, q := range questions {
		question := models.QuizQuestion{
			ID:          uuid.New().String(),
			QuizID:      quiz.ID,
			Type:        models.QuestionMultipleChoice,
			Question:    q.question,
			Explanation: &q.explanation,
			Points:      1,
			SortOrder:   i,
		}
		database.Create(&question)
		for j, opt := range q.options {
			database.Create(&models.QuizOption{
				ID:         uuid.New().String(),
				QuestionID: question.ID,
				Text:       opt.text,
				IsCorrect:  opt.isCorrect,
				SortOrder:  j,
			})
		}
	}

	createSlide(database, s8.ID, "Fleet Management Systems Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Mining Fleet Management Systems — Knowledge Check",
		"inlineQuestions": []map[string]interface{}{
			{
				"question":     "What are the three fundamental functions of every FMS, regardless of vendor or scale?",
				"options":      []string{"Dispatch trucks, record production, send alerts", "Track equipment, measure operational performance, manage site operations", "Optimise haul routes, monitor payload, manage maintenance", "Control autonomous vehicles, report utilisation, manage fuel"},
				"correctIndex": 1,
				"explanation":  "Track (know where equipment is), Measure (accurate production records), Manage (supervisor tools to act on data). Dynamic dispatch is an advanced capability — not a core function present in every FMS.",
			},
			{
				"question":     "A new operation moves from radio dispatch to GPS tracking without changing any other process. Which FMS value ladder step does this represent, and what productivity gain is typical?",
				"options":      []string{"Step 4 — Optimise (+5–8%): deploying the dispatch algorithm first delivers the highest gain", "Step 1 — Track (+5–8%): GPS tracking and cycle counting makes existing inefficiencies visible", "Step 3 — Manage (+8–10%): real-time supervisor tools have the fastest ROI", "Step 2 — Measure (+5–8%): accurate production records are the most urgent first step"},
				"correctIndex": 1,
				"explanation":  "Step 1 — Track — delivers 5–8% productivity improvement simply from GPS tracking and cycle counting, without changing any process. The gains come from eliminating lost trucks, untracked idle time, and inaccurate shift records.",
			},
			{
				"question":     "A blast engineer reports poor fragmentation in one zone despite design penetration rates being achieved. A drill MWD dataset is available. Which parameter most likely explains the fragmentation failure?",
				"options":      []string{"Penetration rate — faster drilling indicates softer rock that produces poor fragmentation", "Hole deviation — misaligned holes reduce effective burden, producing coarse fragmentation", "Rotary pressure and pulldown force — these reveal rock strength variation that penetration rate alone does not", "Bit type — wrong bit causes poor penetration and inadequate fragmentation energy transfer"},
				"correctIndex": 2,
				"explanation":  "Rotary pressure and pulldown force reveal rock strength variation. Higher values in a zone of the same penetration rate indicate harder, more competent rock requiring a different explosive energy factor. MWD-derived lithology inference is the key FMS contribution to blast design quality.",
			},
			{
				"question":     "A match factor of 0.72 at a loader indicates what condition, and why is it the higher-cost failure mode?",
				"options":      []string{"Over-trucked — trucks are queuing, which is the highest-cost failure mode due to multiplied idle costs", "Under-trucked — the loader is starving, which is the highest-cost failure mode due to loader operating costs of $3,000–8,000/hr (rope shovel class)", "Perfectly balanced — MF of 0.72 is within the acceptable target range", "Road congestion — trucks cannot reach the loader fast enough to maintain a healthy match factor"},
				"correctIndex": 1,
				"explanation":  "MF 0.72 is under-trucked — the loader is starving between truck arrivals. Rope shovel and electric mining shovel class equipment costs $3,000–8,000/hr. Loader idle is the highest-cost dispatch failure. A well-designed dispatch system is biased toward slight over-trucking (MF slightly below 1.0).",
			},
			{
				"question":     "Why does payload variance matter more than low average payload when assessing tyre health risk?",
				"options":      []string{"Average payload determines fuel burn per tonne, making it the primary cost driver", "Variance causes inconsistent dispatch cycle times that confuse the optimisation algorithm", "Overloading events within high variance cause non-linear tyre and structural degradation that average payload does not reveal", "FMS systems can only accurately measure variance, not average payload"},
				"correctIndex": 2,
				"explanation":  "Tyre load-life degrades non-linearly with overloading — the relationship is strongly superlinear. Approximately 30% tyre life reduction per 10% overload is a commonly cited approximation. High variance means frequent overloading events that the average figure never reveals.",
			},
			{
				"question":     "What is the most critical integration point between FMS and grade control, and what goes wrong when it fails?",
				"options":      []string{"Real-time mill feed rate reporting — FMS must control the rate of material to the crusher", "Blast movement vectors applied to grade control polygons before loading begins", "CMMS work order generation for drill equipment after MWD anomalies", "Truck speed compliance in grade control exclusion zones"},
				"correctIndex": 1,
				"explanation":  "The blast movement update — applying post-blast movement vectors to grade control polygons before loading begins. When this fails, trucks haul material classified against the pre-blast position (shifted 5–30m). The result is ore hauled to waste or waste to crusher — the highest-cost ore misclassification mechanism in open-cut mining.",
			},
			{
				"question":     "Delay analysis data quality has degraded over 6 months — operators are logging 80% of delays as 'other'. What is the most likely operational consequence for an operation running automated dispatch?",
				"options":      []string{"The FMS network will become congested with unclassified delay messages", "Dispatch assignment quality degrades as the algorithm loses accurate availability predictions, undermining the foundation that optimisation depends on", "The CMMS integration will stop receiving fault codes from the FMS system", "Payload monitoring accuracy will decrease as the sensor calibration drifts"},
				"correctIndex": 1,
				"explanation":  "The dispatch algorithm optimises against predicted cycle times derived from historical data. Delay code breakdown means the system cannot accurately predict when machines will be available — degrading assignment quality. The operation loses the value of Steps 2 and 3 (Measure, Manage), undermining the foundation that Step 4 (Optimise) depends on.",
			},
			{
				"question":     "A supervisor is consistently overriding automated dispatch assignments. What is the most effective intervention?",
				"options":      []string{"Remove override access from supervisor accounts to enforce algorithm compliance", "Track override rates, train supervisors on dispatch algorithm logic, and make override frequency a visible management metric", "Adjust dispatch parameters to match the supervisor's historical manual allocation preferences", "Override behaviour is expected and normal during early FMS deployment — no action needed for 12 months"},
				"correctIndex": 1,
				"explanation":  "Override rates should be tracked and visible as a management metric. High override rates are typically a training issue — the supervisor does not understand why the algorithm made its decision, so they revert to intuition. Training on dispatch logic and making override rates visible to management reduces this without removing the necessary safety valve.",
			},
			{
				"question":     "A 60-truck fleet has Physical Availability 91%, Utilisation of Available Time 89%, but trucks are averaging 108-minute cycles against a 100-minute design. Shift tonnage is 4% below target. What is the primary constraint?",
				"options":      []string{"Physical Availability — schedule additional maintenance to reduce breakdown frequency", "Utilisation of Available Time — reduce queue time at loaders by adjusting dispatch assignments", "Operating efficiency — trucks are running 8% below design cycle time; investigate cycle time vs design by road segment", "Payload — trucks are underloaded, requiring more cycles to hit target tonnage"},
				"correctIndex": 2,
				"explanation":  "Physical Availability (91%) and Utilisation (89%) are both healthy. The constraint is operating efficiency: trucks are running 8% slower than design cycle time, which compounds across 60 trucks. The FMS action is to compare actual cycle times against design by road segment to identify where the slowdown is occurring.",
			},
			{
				"question":     "A 45-truck mixed-OEM open pit operation with no existing GPS tracking is deploying FMS for the first time. Which approach is appropriate?",
				"options":      []string{"Enterprise platform (MineStar Fleet or Modular DISPATCH) — only enterprise systems support mixed-OEM fleets reliably", "OEM-agnostic lightweight FMS — lower cost, faster implementation, no equipment vendor dependency, appropriate to current capability maturity", "Single-OEM FMS from the dominant equipment manufacturer — standardise the fleet first, then deploy FMS", "No FMS needed at 45 trucks — manual dispatch is sufficient at this scale and adds unnecessary overhead"},
				"correctIndex": 1,
				"explanation":  "This operation is at Step 1 of the value ladder — their priority is visibility and basic production records. An OEM-agnostic lightweight system fits: lower implementation cost, no single-vendor dependency, faster to deploy on a mixed fleet, and appropriately scoped to current capability maturity. Enterprise dispatch platforms are purpose-built for large, mature, complex multi-loader environments.",
			},
		},
	})

	createSlide(database, s8.ID, "Module Complete", "COMPLETION", 1, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "Congratulations — you have completed the Mining Fleet Management Systems module.",
		"keyTakeaways": []string{
			"FMS is a platform to measure, manage, and optimise — not primarily a dispatch tool",
			"The value ladder compounds: Track → Measure → Manage → Optimise delivers up to 30% productivity improvement",
			"Dynamic dispatch is an advanced capability — it requires the Track and Measure foundations to be solid first",
			"Payload consistency matters more than average payload — tyre degradation is non-linear with overload",
			"Data quality is the highest-leverage investment in any FMS programme — delay codes, payload capture, and auto-classification rate are the KPIs that determine whether everything else works",
			"Most FMS underperformance is organisational, not technical — supervisor training and data discipline are the interventions that move the needle",
		},
		"showCertificate": true,
	})

	log.Println("Fleet Management Systems module seeded with 8 sections, diagrams, day-in-the-life slides, and quiz")
}
