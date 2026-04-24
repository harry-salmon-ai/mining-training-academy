package seed

import (
	"encoding/json"
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedCollisionAvoidanceModule(database *gorm.DB, force bool) {
	var existing models.Module
	database.Where("slug = ?", "collision-avoidance-systems").First(&existing)
	if existing.ID != "" {
		if force {
			log.Println("SEED=true: deleting existing Collision Avoidance Systems module for re-seed...")
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
			log.Println("Existing Collision Avoidance Systems module deleted")
		} else {
			log.Println("Collision Avoidance Systems module already exists, skipping")
			return
		}
	}

	log.Println("Seeding Collision Avoidance Systems module...")

	var category models.Category
	database.Where("slug = ?", "safety-compliance").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed CAS: category or admin not found")
		return
	}

	desc := "Collision Avoidance Systems — The Perception-Based Approach. Why the industry says it can't be built, and why it can. Covers the EMESRT vehicle interaction framework, NSW TRG ROVOA regulatory requirements, the structural failure modes of transponder-based systems — including the false positive taxonomy and the operational memory gap — and the perception-based architecture that addresses them. Designed for engineers and safety practitioners building a rigorous, technically grounded understanding of CAS."
	dur := 90
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Collision Avoidance Systems — The Perception-Based Approach",
		Slug:        "collision-avoidance-systems",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelAdvanced,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"CAS", "collision avoidance", "EMESRT", "NSW TRG", "ROVOA", "perception", "4D radar", "LiDAR", "vehicle interaction", "safety", "alarm fatigue", "Level 9", "PR-5A", "transponder", "false positives", "operational memory"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 0: The Vehicle Interaction Problem ──────────────────────────────
	s0 := createSection(database, module.ID, "The Vehicle Interaction Problem", 0)

	casContentJSON, _ := json.Marshal(map[string]interface{}{
		"type":     "title",
		"heading":  "Collision Avoidance Systems — The Perception-Based Approach",
		"subtitle": "The technology the industry says can't be built — and why it can",
		"bullets": []string{
			"Approximately 25–35% of annual mining fatalities are attributable to vehicle interaction failures (ICMM Health & Safety Performance data)",
			"NSW Resources Regulator has identified vehicle/pedestrian interaction as the #1 hazard in open cut mining",
			"Regulatory requirements have accelerated: NSW TRG ROVOA (July 2025) mandates CAT Layers 8 and 9",
			"Perception-based systems detect any physical object regardless of electronic tagging — the architectural shift the industry needs",
		},
	})
	database.Create(&models.Slide{ID: uuid.New().String(), Title: "Collision Avoidance Systems — The Perception-Based Approach", Type: "TITLE", Content: casContentJSON, SortOrder: 0, SectionID: s0.ID})

	createSlide(database, s0.ID, "The #1 Hazard in Open Cut Mining", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The #1 Hazard in Open Cut Mining",
		"introText": "Vehicle interaction is not an emerging risk — it is the most persistent, most costly, and most analysed safety problem in open cut mining. Despite decades of investment, the fatality rate has not moved.",
		"items": []map[string]interface{}{
			{
				"title":       "The scale of the problem",
				"description": "ICMM Health & Safety Performance data consistently places vehicle interaction failures at approximately 25–35% of annual mining fatalities. On a global industry operating at scale, this translates to dozens of preventable deaths every year.",
			},
			{
				"title":       "The NSW Regulator's position",
				"description": "The NSW Resources Regulator has formally identified vehicle/pedestrian interaction as the #1 hazard in open cut mining — a designation that drives enforcement posture, mandatory standards, and the July 2025 Technical Reference Guide (ROVOA).",
			},
			{
				"title":       "Why it persists despite investment",
				"description": "Tier 1 miners have invested hundreds of millions in CAS deployments since 2015. The fatality rate has not materially improved. The problem is structural, not technological — existing technology architectures have fundamental limitations that money alone cannot solve.",
			},
			{
				"title":       "The structural gap",
				"description": "Current controls rely on either tagged equipment (every hazard must carry an active transponder) or operator vigilance (a tired operator in a 220-tonne haul truck must notice a pedestrian). Neither is reliable at scale. The solution requires a fundamentally different sensing architecture.",
			},
		},
		"keyFact": "Approximately 25–35% of mining fatalities each year are attributable to vehicle interaction failures — ICMM Health & Safety Performance data.",
	})

	createSlide(database, s0.ID, "Why Existing Controls Keep Failing", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Why Existing Controls Keep Failing",
		"introText": "Three structural failure modes explain why CAS deployments have not moved the fatality needle. Each is addressable — but not by incremental improvement of existing architectures.",
		"items": []map[string]interface{}{
			{
				"title":       "Alarm fatigue destroys operator trust",
				"description": "The Glencore/Wabtec Gen 1 deployment recorded 50,000 alarms on just two dozers in a single operating period. When a system is wrong 99% of the time, operators rationally stop responding — including to the 1% that are genuine threats. Alarm fatigue is not an operator failure; it is a system design failure.",
			},
			{
				"title":       "Transponder dependency creates dangerous blind spots",
				"description": "Tag-based systems are invisible to any untagged hazard. Contractors arriving without site tag configuration, maintenance crews whose PPE tag battery is dead, rock falls, mis-parked equipment, and third-party vehicles are all completely invisible to a transponder CAS. The assumption of 100% tag compliance on a large site is operationally unsustainable.",
			},
			{
				"title":       "Compliance does not equal capability",
				"description": "A system can be fully PR-5A compliant — meeting every EMESRT performance requirement — and still fail operationally. PR-5A defines detection performance in controlled test scenarios. It does not define alarm quality, contextual intelligence, or real-world NVA rates in a working mine.",
			},
			{
				"title":       "The regulatory gap",
				"description": "Regulators have been clear about what they need: a system that detects any physical hazard and can intervene when the operator does not respond. What vendors have delivered — until now — is a system that detects tagged hazards and generates alarms. The gap between regulatory need and vendor delivery is structural, not incremental.",
			},
		},
	})

	// ── Section 1: The Vehicle Interaction Control Framework (CONSOLIDATED) ─────
	s1 := createSection(database, module.ID, "The Vehicle Interaction Control Framework", 1)

	createSlide(database, s1.ID, "The Vehicle Interaction Control Framework", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "The Vehicle Interaction Control Framework",
		"subtitle": "EMESRT — the industry standard, NSW TRG — the regulatory codification",
		"bullets": []string{
			"9-level framework spanning Design (L1–3), Operate (L4–6), and React (L7–9) tiers",
			"Design controls must precede Operate controls — which must precede React (technology) controls",
			"PR-5A performance requirements define what a CAS must demonstrate",
			"NSW TRG ROVOA (July 2025) codifies the same framework into NSW regulatory requirements",
		},
	})

	createSlide(database, s1.ID, "EMESRT Vehicle Interaction Control — 9 Levels", "PROCESS", 1, map[string]interface{}{
		"type":        "process",
		"heading":     "EMESRT Vehicle Interaction Control — 9 Levels",
		"description": "The EMESRT framework defines 9 progressive levels of vehicle interaction control. Levels 1–3 are Design controls — built into the physical and procedural fabric of the mine. Levels 4–6 are Operate controls. Levels 7–9 are React (technology-assisted) controls — the domain where CAS products operate. Design before Operate before React is not a preference; it is the fundamental sequencing requirement.",
		"steps": []map[string]interface{}{
			{
				"title":       "Level 1 — Site Requirements",
				"tier":        "Design",
				"description": "Equipment specification, standards, mine designs and plans. The foundation for all subsequent controls — decisions made here shape the hazard environment for the life of the mine.",
			},
			{
				"title":       "Level 2 — Segregation Controls",
				"tier":        "Design",
				"description": "Physical barriers, control gates, road segregation, and activity scheduling. Structural separation of vehicles and pedestrians wherever possible before operational procedures are required.",
			},
			{
				"title":       "Level 3 — Physical & Procedural Design Controls",
				"tier":        "Design",
				"description": "Physical engineering controls (barriers, bollards, traffic islands, one-way ramps) and procedures designed into the mine's operating methodology before operations begin. Engineering controls that do not rely on operator compliance or technology function.",
			},
			{
				"title":       "Level 4 — Authority to Operate",
				"tier":        "Operate",
				"description": "Vehicle-specific training, licensing, and operating area access control. Ensures only trained and authorised personnel operate equipment in defined areas.",
			},
			{
				"title":       "Level 5 — Fitness to Operate",
				"tier":        "Operate",
				"description": "Physical, psychological, fatigue, and substances assessment. Addresses the human impairment factors that degrade operator vigilance and decision-making.",
			},
			{
				"title":       "Level 6 — Operating Compliance",
				"tier":        "Operate",
				"description": "Vehicle pre-start, defect recording and reporting, traffic management plan adherence. Ensures equipment in service meets minimum vehicle standards and that operating procedures are followed.",
			},
			{
				"title":       "Level 7 — Operator Awareness",
				"tier":        "React",
				"description": "Cameras, proximity maps, real-time alerts, lighting, and mirrors. The technology layer providing the operator with situational awareness beyond direct line of sight. Entry point for CAS technology investment.",
			},
			{
				"title":       "Level 8 — Advisory Controls",
				"tier":        "React",
				"description": "Proximity alerts, overspeed warnings, and fatigue monitoring. Alerts that require the operator to take action — the system advises, the operator decides and responds.",
			},
			{
				"title":       "Level 9 — Machine Intervention",
				"tier":        "React",
				"description": "Automatic braking, speed control, and motion inhibit. The machine takes action without waiting for the operator when a confirmed collision threat is detected.",
				"warnings": []string{
					"Level 9 depends entirely on Level 7 and Level 8 data quality. A Level 9 system operating on a high-NVA advisory layer triggers unsafe interventions on non-genuine threats. Do not activate intervention before advisory performance is validated.",
				},
			},
		},
	})

	createSlide(database, s1.ID, "NSW TRG ROVOA — The Regulatory Codification", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "NSW TRG ROVOA — The Regulatory Codification",
		"introText": "The NSW Resources Regulator's Technical Reference Guide (July 2025) codifies the EMESRT framework into NSW regulatory requirements. The structure is the same — 9 layers, Design before Operate before React — with specific enforcement implications.",
		"items": []map[string]interface{}{
			{
				"title":       "Same framework, regulatory force",
				"description": "The NSW TRG maps directly onto the EMESRT 9-level framework. Layers 1–3 are Design controls. Layers 4–7 are Operate controls. Layers 8–9 are CAT (Collision Avoidance Technology) — the technology intervention layers. The Regulator's designation of these as mandatory layers, not aspirational guidance, is the critical difference.",
			},
			{
				"title":       "The enforcement position",
				"description": "The NSW Chief Inspector of Mines: \"Technology is part of the solution, but it is not the solution. If you don't get the fundamentals right — design, layout, vehicle specifications — no amount of technology will compensate.\" 144 unannounced site inspections in 2024–25 issued prohibition notices where fundamental controls were absent regardless of technology deployment.",
			},
			{
				"title":       "The path to legislative codification",
				"description": "The TRG ROVOA is the pre-legislative signal. The Regulator has indicated these requirements will be codified in the Mining and Petroleum Legislation Amendment Regulation. Compliance with the TRG now is the pathway to compliance with future legislation. Sites that deploy technology without Layers 1–7 are setting themselves up for prohibition notices today and legislative non-compliance tomorrow.",
			},
			{
				"title":       "Implication for CAS vendors",
				"description": "A technology vendor that sells Layer 8/9 capability without helping clients implement Layers 1–7 is setting its clients up for regulatory non-compliance and operational failure. The strongest CAS deployments include a Layers 1–7 gap analysis as part of the scope of work.",
			},
		},
	})

	createSlide(database, s1.ID, "CAT Layer 8 and 9 — What They Actually Demand", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "CAT Layer 8 and 9 — What They Actually Demand",
		"description": "The NSW TRG and EMESRT PR-5A define specific performance requirements for CAT Layers 8 and 9. These are not aspirational targets — they are the baseline that a compliant system must demonstrate in production operations.",
		"diagramType": "kpi-map",
		"diagramData": map[string]interface{}{
			"ratios": []map[string]interface{}{
				{
					"name":        "CAT Layer 8 — Warning Detection Range",
					"formula":     "Heavy vehicle >30m | Light vehicle >20m | Person >10m at walking speed",
					"owner":       "Technology / Operations",
					"description": "Minimum detection ranges for each hazard class at Zone 2 entry. These distances are calibrated to operator reaction time and vehicle braking distance at typical haul road speeds. Detection below these thresholds leaves insufficient time for the operator to respond before Zone 1 entry.",
				},
				{
					"name":        "CAT Layer 8 — Alert Specificity",
					"formula":     "NVA rate <2 per machine-hour (heavy fleet)",
					"owner":       "Technology / Operations",
					"description": "False positive rate must remain below operational tolerance. The ratio of value-add alerts to total alerts determines whether operators trust and respond to the system. A high NVA rate produces alarm fatigue — the primary failure mode of Gen 1 CAS deployments. Target: <2 NVA alerts per machine-hour in production.",
				},
				{
					"name":        "CAT Layer 9 — Intervention Latency",
					"formula":     "Detection-to-intervention <500ms",
					"owner":       "Technology",
					"description": "End-to-end time from confirmed threat detection to speed reduction initiation. At 40 km/h relative closing speed, 500ms equals ~5.5m of approach distance — a non-negotiable budget on high-speed haul segments. Latency above 500ms means the vehicle enters Zone 1 before the system can respond.",
				},
				{
					"name":        "CAT Layer 9 — Intervention Reliability",
					"formula":     "System availability >99.5% | Failsafe: controlled stop on fault",
					"owner":       "Technology",
					"description": "System availability during scheduled operating time and failsafe behaviour on system fault. Safe state on any sensor or compute fault must be a controlled deceleration to safe speed — not loss of control, and not silent failure. The system must self-diagnose and alert to degraded mode.",
				},
			},
		},
	})

	createSlide(database, s1.ID, "PR-5A — Performance Requirements for Vehicle Interaction Systems", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "PR-5A — Performance Requirements for Vehicle Interaction Systems",
		"introText": "EMESRT PR-5A is the specification that defines what a CAS must demonstrate. It is the baseline standard for procurement decisions across ICMM member companies and the reference for all compliance claims.",
		"items": []map[string]interface{}{
			{
				"title":       "24 surface vehicle interaction scenarios",
				"description": "PR-5A defines 24 discrete scenarios covering vehicle-to-vehicle, vehicle-to-pedestrian, and vehicle-to-static hazard interactions. Each scenario has defined approach geometry, speeds, and minimum performance requirements. A compliant system must pass all 24 — including oblique and rear approaches that are harder to detect than head-on encounters.",
			},
			{
				"title":       "Alert escalation and system integrity",
				"description": "PR-5A defines alert escalation logic (awareness → advisory → intervention), self-diagnostic requirements, and failover behaviour. A compliant system must report its own health status and behave safely on internal fault.",
			},
			{
				"title":       "The critical limitation of compliance alone",
				"description": "PR-5A compliance is measured in controlled test scenarios against tagged hazards. It does not measure NVA alarm rate in production operations, operator acceptance, or alarm fatigue risk. The Glencore Gen 1 deployment was PR-5A compliant — and generated 50,000 alerts on two dozers. Compliance is the floor, not the ceiling.",
			},
			{
				"title":       "The key advantage of perception over transponder for all 24 scenarios",
				"description": "A transponder system handles all 24 scenarios — if and only if every hazard carries an active, functioning electronic tag. Perception handles all 24 scenarios regardless of whether any party carries a tag. This is the architectural difference with direct safety implications.",
			},
		},
		"keyFact": "PR-5A compliance is the minimum standard for procurement — but operational success requires NVA rate management that PR-5A does not measure.",
	})

	// ── Section 2: Why Existing CAS Solutions Fall Short ───────────────────────
	s2 := createSection(database, module.ID, "Why Existing CAS Solutions Fall Short", 2)

	createSlide(database, s2.ID, "The Gap in Current Solutions", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "The Gap in Current Solutions",
		"subtitle": "Why the industry invested billions and the fatality rate hasn't moved",
		"bullets": []string{
			"Transponder systems assume a tagged world — reality is never fully tagged",
			"Alarm fatigue is the #1 deployment failure mode — caused by false positives, not sensor failure",
			"False positives have three structural sources — none of which are solved by better sensors alone",
			"Operational memory is the missing sixth control — no system has it yet",
		},
	})

	createSlide(database, s2.ID, "The Transponder Problem — A Tagged World Assumption", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Transponder Problem — A Tagged World Assumption",
		"introText": "Every transponder-based CAS is built on the same foundational assumption: every hazard carries an active, functioning electronic tag. On a working mine site, this assumption fails continuously.",
		"items": []map[string]interface{}{
			{
				"title":       "Untagged hazards are completely invisible",
				"description": "A transponder CAS has no mechanism to detect an object that does not emit a signal. Contractor arriving on site without tag configured, maintenance crew working in the pit without PPE tag online, visitor vehicle from an adjacent property — all are invisible. The system does not know what it cannot see.",
			},
			{
				"title":       "Tag failure modes are operationally common",
				"description": "Dead battery, forgotten PPE (hard hat left in the crib room), signal interference from metal structures, and tag hardware failure all create invisible hazards. On a site with 500 tagged personnel and vehicles, achieving 100% tag compliance on every shift is a full-time compliance management task.",
			},
			{
				"title":       "The maintenance overhead compounds at scale",
				"description": "Every vehicle and every pedestrian on a large site must be individually tagged, tags must be maintained and batteries replaced, and compliance must be audited. For a 2,000-person operation with 300 vehicles, this is a significant operational overhead — and any gap creates an invisible hazard.",
			},
			{
				"title":       "Perception removes the tagging requirement entirely",
				"description": "A camera/radar/LiDAR sees what is physically present, not what the network knows about. An untagged contractor vehicle is detected the same way a tagged haul truck is detected — by the physical fact of its presence in the detection zone. The architectural shift eliminates the entire category of tag-dependent blind spots.",
			},
		},
	})

	createSlide(database, s2.ID, "Alarm Fatigue — The Silent Killer of CAS Programmes", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Alarm Fatigue — The Silent Killer of CAS Programmes",
		"introText": "Alarm fatigue is not an operator discipline problem. It is a predictable, rational response to a system that is consistently wrong. When a CAS generates overwhelming volumes of non-value-add alarms, operators adapt — and the adaptation is to stop responding.",
		"items": []map[string]interface{}{
			{
				"title":       "The Glencore/Wabtec Gen 1 case",
				"description": "50,000 alerts recorded on two dozers in a single operating period. This figure — drawn from a real deployment — represents the consequence of detecting proximity without context. Every vehicle passing on a separated haul road, every nearby machine at a dump, every slow-moving support vehicle — all triggered the same alert as an actual collision threat.",
			},
			{
				"title":       "The rational operator response",
				"description": "An operator who has received 200 alarms in a shift with no near-miss events learns that the system is noise. When the system is wrong 99% of the time, the correct Bayesian response is to discount all alarms — including the 1% that are genuine. Alarm fatigue is operators behaving logically in response to a poorly calibrated system.",
			},
			{
				"title":       "Why a high-NVA system is worse than no system",
				"description": "A system that generates overwhelming NVA alarms does not merely fail to improve safety — it actively degrades it. Operators who have learned to ignore the system are less alert than operators who have no system and rely on their own vigilance. A broken CAS erodes the safety culture it was meant to support.",
			},
			{
				"title":       "The fix is contextual intelligence, not better sensors",
				"description": "The Glencore Gen 1 sensors were detecting correctly. The problem was the absence of contextual classification: the system could not distinguish 'a vehicle on a parallel road 25 metres away travelling in the same direction' from 'a vehicle on a collision heading 25 metres away'. Perception-based ML classification provides that distinction — but only where it has the context to do so.",
			},
		},
		"keyFact": "50,000 alerts on two dozers (Glencore/Wabtec, conference-reported operational data) — alarm fatigue, not sensor failure, is the primary CAS deployment risk.",
	})

	createSlide(database, s2.ID, "The False Positive Taxonomy — Three Structural Sources", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "The False Positive Taxonomy — Three Structural Sources",
		"introText": "False positives in CAS deployments are not random noise — they have specific, predictable causes. Understanding the taxonomy is essential because each source requires a different solution. Better sensors solve none of them; contextual intelligence is required for all three.",
		"items": []map[string]interface{}{
			{
				"title":       "1. Proximity without path context",
				"description": "The simplest and most common false positive: a vehicle enters the detection zone but its trajectory is diverging, not converging. A parallel haul road 25m away generates identical zone-entry events to an oncoming collision course. Rule-based systems cannot distinguish them without path prediction. This is the false positive the Glencore Gen 1 deployment suffered at scale.",
			},
			{
				"title":       "2. Path projection without mine geometry — the switchback problem",
				"description": "At the apex of a switchback or hairpin ramp turn, a haul truck's current kinematic vector points directly at the oncoming lane. A system that projects linear paths from current velocity flags this as a collision course — every single time. The truck is executing a completely normal, road-geometry-constrained manoeuvre. Without a model of mine road geometry, the system has no way to distinguish this from an unexpected heading change. This fires on every switchback, every shift, at every site with hairpin ramps. It is not an edge case — it is structural.",
			},
			{
				"title":       "3. Positioning without mine geometry — the multi-bench problem",
				"description": "GNSS-based collision detection evaluates risk in 2D or with poor vertical resolution. A dozer working on bench floor 1 and a truck operating on bench 3 directly above appear co-located in 2D coordinates. A pit wall and 20-30 metres of vertical separation mean they cannot physically interact — but the positioning system has no model of pit walls, bench faces, or bench height geometry, so it flags a collision event. Perception-based systems cannot see through pit walls — the physical world filters this false positive automatically. Positioning systems require explicit mine geometry context to achieve the same result.",
			},
		},
		"keyFact": "False positives are not a sensor performance problem — they are a context problem. A system that detects correctly but has no model of the mine's geometry or operational patterns will generate structural false positives regardless of sensor quality.",
	})

	createSlide(database, s2.ID, "Operational Memory — The Missing Sixth Control", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Operational Memory — The Missing Sixth Control",
		"introText": "Every CAS failure mode discussed so far has a named cause. There is a sixth structural gap that has no name in the industry literature — but explains why false positive rates don't improve over time on any deployed system. Every event is evaluated in isolation. The system has no memory.",
		"items": []map[string]interface{}{
			{
				"title":       "No model of what normal looks like on this mine",
				"description": "A truck turning at the top of a ramp at the end of a haul cycle generates identical alert signals to a truck unexpectedly changing course toward a pedestrian work area — because the system has no knowledge of where trucks normally go. The mine's operational envelope — where machines travel, when they turn, what approach patterns precede a dump — is not encoded anywhere. The system never learns it.",
			},
			{
				"title":       "Why false positive rates don't improve over time",
				"description": "Without an accumulating operational context, a deployed CAS has no baseline to improve against. It generates the same switchback false positive on day 500 as on day 1. There is no feedback loop from observed machine behaviour to prediction model. This is why alarm fatigue persists in long-running deployments — the system cannot get smarter about this mine.",
			},
			{
				"title":       "What operational memory would change",
				"description": "A system that builds its own model of the mine's operational patterns — where machines normally travel, when they turn, what constitutes an anomalous path — from observed machine movement changes the character of the problem. Deviation from the learned operational baseline becomes the trigger, not proximity alone. 'This truck is heading somewhere trucks do not normally go' is a far more actionable signal than 'this truck is within 25m of another truck'.",
			},
			{
				"title":       "The continuous improvement loop",
				"description": "A system with operational memory gets smarter over time without external integration. It does not need an FMS feed to know there is a switchback at this grid location — it has observed thousands of trucks navigate it correctly. Each shift's data refines the model. This is the architectural property that makes false positive rates trend toward zero over a deployment lifetime, rather than remaining structurally fixed.",
			},
		},
		"keyFact": "Operational memory is not in any deployed CAS today. It is the next frontier — and the property that separates systems that get smarter from systems that stay broken.",
	})

	createSlide(database, s2.ID, "The Hardware Constraints That Remain", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "The Hardware Constraints That Remain",
		"introText": "Two of the original objections to perception-based CAS are genuine, persistent hardware limitations — not problems that have been solved. Honest framing matters: these are managed, not eliminated.",
		"items": []map[string]interface{}{
			{
				"title":       "Dust — managed, not solved",
				"description": "Open cut mines generate dust clouds that attenuate optical systems. Cameras and LiDAR lose detection performance in heavy mineral dust — this is a real and acknowledged limitation. It is not a limitation of radar. The 4D radar sensor operates on radio waves, which are not materially attenuated by particulate matter. Sensor fusion handles dust by switching to radar-primary mode, with graceful degradation of classification fidelity rather than system failure. But dust does degrade the perception system — calling it 'solved' overstates the engineering reality.",
			},
			{
				"title":       "Range — managed at the haul road level, constrained elsewhere",
				"description": "Automotive AI cameras have reliable detection range of approximately 15m under mining conditions — below the 30m+ required for haul truck deceleration at speed. This is addressed by medium-range telephoto lens integration (30–80m reliable range) and radar's 100m+ range. For haul road applications, the range problem is adequately addressed. For pedestrian detection in confined areas with dust, range constraints remain a genuine operational limitation that site-specific validation must quantify.",
			},
		},
	})

	createSlide(database, s2.ID, "The Barriers That Have Fallen", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "The Barriers That Have Fallen — The Enablers That Arrived",
		"introText": "Two of the five original objections to perception-based CAS are no longer blockers. Hardware and methodology that did not exist at production-grade price and performance points before 2022 now does. This is why the timing argument matters — perception-based CAS was not buildable five years ago.",
		"items": []map[string]interface{}{
			{
				"title":       "Edge compute — resolved by NVIDIA Orin and equivalent SoCs (2022+)",
				"description": "Real-time 3D perception processing at sub-200ms latency on a vibrating haul truck required automotive-grade hardware with dedicated neural accelerators. Until 2022–2023 (NVIDIA Orin, Qualcomm Snapdragon Ride), this hardware did not exist at industrial price and power points. NVIDIA Orin delivers >100 TOPS at 15–65W — real-time 3D inference on a haul truck is now achievable within mining power and thermal budgets. This is a solved problem.",
			},
			{
				"title":       "Training data — resolved by synthetic data + transfer learning",
				"description": "Distinguishing a person walking in dust from a low-visibility background object required training data from mining environments. Automotive datasets do not transfer directly. Mining-specific training data did not exist at scale until dedicated collection programmes started in 2021–2023. Photorealistic simulation now generates ground-truth labelled training data for edge cases that are rare or dangerous to capture in production. Transfer learning from automotive datasets (millions of labelled instances) plus 6–12 months of targeted mining-environment collection closes the remaining gap.",
			},
		},
		"keyFact": "Edge compute and training data are no longer reasons perception-based CAS cannot be built. They are engineering work that is underway — the barriers have fallen.",
	})

	// ── Section 3: The Perception-Based Approach ────────────────────────────────
	s3 := createSection(database, module.ID, "The Perception-Based Approach — Why Now", 3)

	createSlide(database, s3.ID, "The Perception-Based Approach", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "The Perception-Based Approach",
		"subtitle": "Why the silver bullet can now be built",
		"bullets": []string{
			"Detect, classify, and track any physical object — regardless of electronic tagging",
			"4D imaging radar + AI cameras + LiDAR: all-conditions coverage no single sensor provides",
			"Automotive-grade edge compute (NVIDIA Orin) makes real-time 3D inference viable at mining price points",
			"Contextual ML classification reduces NVA rates — operational memory will reduce them further",
		},
	})

	createSlide(database, s3.ID, "What Perception-Based Means", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What Perception-Based Means",
		"introText": "A perception-based CAS detects, classifies, and tracks any physical object in the detection zone — regardless of whether that object carries any electronic tag. The sensor sees what is physically present, not what the network knows about.",
		"items": []map[string]interface{}{
			{
				"title":       "The architectural definition",
				"description": "Camera, LiDAR, and radar sensors observe the physical environment directly. Object detection and classification algorithms run on onboard edge hardware in real time. The system's knowledge of what is in the scene comes from sensor data, not from a communications network of transponders.",
			},
			{
				"title":       "The passenger vehicle parallel",
				"description": "This is the same architectural shift that made Tesla Autopilot, AEB (Automatic Emergency Braking), and lane-departure warning standard features in modern passenger vehicles. Those systems do not require other vehicles or pedestrians to carry transponders — they see the physical world. The mining environment is harder, but the fundamental approach is identical.",
			},
			{
				"title":       "Why the mining environment is harder",
				"description": "Dust, vibration, extreme temperature ranges, GNSS denial in deep pits, 24-hour operation, and the scale of vehicles involved (a 220-tonne haul truck has different braking dynamics from a passenger car) all add complexity. These are engineering problems — they do not invalidate the architectural approach, they define the specification for solving it.",
			},
			{
				"title":       "4D radar + AI cameras + LiDAR = all-conditions coverage",
				"description": "No single sensor provides all-conditions coverage. Cameras provide rich classification data in good visibility; radar provides all-weather detection and velocity measurement regardless of dust or fog; LiDAR provides precise 3D ranging in normal conditions. Sensor fusion means the system degrades gracefully in adverse conditions rather than failing catastrophically.",
			},
		},
	})

	createSlide(database, s3.ID, "Sensor Fusion Stack", "DIAGRAM", 2, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Sensor Fusion Stack — The Perception Architecture",
		"description": "A perception-based CAS fuses data from four sensor modalities through a five-stage core algorithm. Each sensor contributes different capabilities; the fusion layer combines them into a single, coherent threat picture.",
		"diagramType": "architecture",
		"diagramData": map[string]interface{}{
			"sensors": []map[string]interface{}{
				{
					"name":        "AI Vision Cameras",
					"description": "Wide-angle and telephoto cameras with onboard inference. Detect, classify, and track vehicles, persons, and static hazards with >95% classification accuracy in normal conditions. Telephoto lens extends reliable range to 30–80m for haul road applications.",
					"functions":   []string{"Object classification", "Person detection", "Range estimation", "Visual context"},
				},
				{
					"name":        "4D Imaging Radar",
					"description": "FMCW radar with Doppler velocity measurement. All-weather performance unaffected by dust, rain, or fog. Operates in GNSS-denied pit environments. Measures both position AND velocity — enabling collision prediction rather than just proximity detection.",
					"functions":   []string{"All-weather detection", "Velocity measurement", "GNSS-denied zones", "V2V ranging"},
				},
				{
					"name":        "Solid-State LiDAR",
					"description": "360° or sector-scan 3D point cloud. Sub-centimetre resolution at range in normal conditions. Performance degrades in heavy dust but less severely than cameras. Provides the depth map that anchors trajectory prediction and path geometry.",
					"functions":   []string{"3D point cloud", "Precise ranging", "Terrain mapping", "Void detection"},
				},
				{
					"name":        "GNSS / RTK Positioning",
					"description": "High-precision vehicle position where satellite reception is available. Combined with dead-reckoning via IMU for GNSS-shadowed zones (deep pits, under berms). Not the primary perception layer — provides geo-context for zone enforcement and fleet-level situational awareness.",
					"functions":   []string{"Vehicle position", "Map matching", "Zone enforcement", "Fleet context"},
				},
			},
			"coreAlgorithm": map[string]interface{}{
				"stages": []string{
					"Detect\n(multi-modal fusion)",
					"Classify\n(ML inference)",
					"Track\n(trajectory predict)",
					"Threat assess\n(context engine)",
					"Alert / Intervene\n(L7/8/9 output)",
				},
			},
		},
	})

	createSlide(database, s3.ID, "CAS Sensors on a Haul Truck", "IMAGE", 3, map[string]interface{}{
		"type":        "image",
		"heading":     "CAS Sensors on a Haul Truck",
		"description": "A perception-based CAS is a physical installation on each vehicle. The radar antennas, camera housings, LiDAR units, and GNSS antenna are mounted at specific locations to achieve full-perimeter coverage. Unlike a transponder badge clipped to a belt, the sensor suite is engineered infrastructure — always on, no batteries to change, no operator action required.",
		"imageUrl":    "https://assets.miningacademy.com/images/cas-sensors-haul-truck-mounted.jpg",
	})

	createSlide(database, s3.ID, "The Three Detection Zones", "DIAGRAM", 4, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Three Detection Zones",
		"description": "The zone model translates EMESRT Levels 7, 8, and 9 into spatial detection boundaries around the host vehicle. Zone sizes are not static — they scale dynamically with closing speed and braking distance.",
		"diagramType": "layers",
		"diagramData": map[string]interface{}{
			"layers": []map[string]interface{}{
				{
					"label":       "Zone 3 — Awareness / Information (>30m)",
					"description": "Operator awareness: nearby equipment visible on display, context information provided. No alarm generated. Equivalent to EMESRT Level 7 — the operator is informed, not alerted. Berms, pit walls, and separated-road vehicles appear in Zone 3 and do not escalate.",
					"stage":       "Perception",
				},
				{
					"label":       "Zone 2 — Warning / Advisory (15–30m)",
					"description": "Audible and visual warning alert. Operator is expected to assess the hazard and take evasive action. System monitors operator response. Equivalent to EMESRT Level 8. Alert must be specific to object type, heading, and closing speed — not a generic proximity alarm.",
					"stage":       "Comprehension → Projection",
				},
				{
					"label":       "Zone 1 — Emergency / Intervention (<15m)",
					"description": "Imminent collision. If the operator has not responded to Zone 2 advisory, the machine reduces speed and applies controlled braking. Equivalent to EMESRT Level 9. Zone 1 boundary expands dynamically at high relative speeds — at 40 km/h combined closing speed, braking distance alone requires the Zone 1 threshold to extend well beyond 15m.",
					"stage":       "Decision → Action",
				},
			},
			"contextFactors": []string{
				"Closing speed",
				"Object classification",
				"Heading geometry",
				"Road grade",
			},
			"callout": map[string]interface{}{
				"text":        "Zone distances are indicative baselines — real deployments use dynamic zone sizing based on closing speed and braking distance.",
				"description": "Zone 1 <15m and Zone 2 15–30m are baseline values for a heavy vehicle at typical haul road speed. At higher closing speeds, or for pedestrian detection where braking distances are shorter, zone boundaries are recalculated in real time. A static zone boundary fails above ~25 km/h combined closing speed.",
			},
		},
	})

	createSlide(database, s3.ID, "EMESRT Vehicle Interaction Zones", "DIAGRAM", 5, map[string]interface{}{
		"type":        "diagram",
		"heading":     "EMESRT 3-Zone Vehicle Interaction Model",
		"description": "The EMESRT 3-zone model — Awareness (Zone 3), Warning (Zone 2), Intervention (Zone 1) — in plan view around a heavy vehicle. Each concentric zone corresponds to an EMESRT control level (L7, L8, L9) and drives a different system response: display-only in Zone 3, operator alert in Zone 2, automatic intervention in Zone 1. Zone boundaries are not fixed circles — they stretch dynamically with closing speed and braking distance, and asymmetrically to reflect sensor coverage and blind-spot geometry around the host vehicle.",
		"diagramType": "image",
		"imageUrl":    "/assets/diagrams/cas_vehicle_interaction_zones.png",
	})

	createSlide(database, s3.ID, "Sensor Modalities — Tradeoffs for CAS", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Sensor Modalities — Tradeoffs for CAS",
		"introText": "No single sensor covers every failure mode. Production CAS systems fuse multiple modalities because each has blind spots that the others fill. Understanding the tradeoffs is fundamental to designing the sensor suite, power budget, and compute architecture.",
		"items": []map[string]interface{}{
			{"title": "4D Radar", "description": "Long range (100m+), works in dust, rain, and darkness, directly measures velocity via Doppler. Tradeoff: low angular resolution (a few degrees), poor object classification — a radar return can tell you something is there and moving, but not easily whether it's a truck, a light vehicle, or a berm."},
			{"title": "LiDAR", "description": "High-resolution 3D point cloud, excellent for object shape classification and precise localisation. Tradeoff: degrades in heavy dust and rain (returns are scattered and attenuated), requires significantly more compute for point cloud processing, higher cost per unit, mechanical scanning LiDARs have moving parts that fail in harsh mine environments."},
			{"title": "Camera / Vision", "description": "Richest per-pixel data, enables detailed object classification (truck vs person vs dozer), reads signage and paint lines. Tradeoff: lighting-dependent (fails in dust, fog, direct sun glare, and at night without IR), heavy compute for real-time CNN inference, high false-positive risk on dust particles and shadows unless carefully tuned."},
			{"title": "Ultrasonic", "description": "Short-range (1–5 m), simple, cheap, extremely reliable. Used for low-speed proximity detection — parking sensors on passenger cars, final-metre confirmation on heavy equipment. Tradeoff: cannot substitute for longer-range sensors; only useful in the innermost detection envelope."},
			{"title": "Fusion Rationale", "description": "No modality is sufficient alone. Production systems run radar + LiDAR as the primary fusion pair — radar gives range and velocity reliably in bad weather, LiDAR gives resolution and classification. Camera is supplemental for classification confirmation and operator display. Ultrasonic handles the final metre."},
			{"title": "Why This Matters for Engineers", "description": "Sensor selection cascades into every other design decision: mounting positions (radar needs clear forward view, LiDAR needs roof line, camera needs wiper-protected glass), power budget (LiDAR and camera compute drive >100 W continuous), compute architecture (real-time fusion pipeline topology), and failure modes (what happens when LiDAR is obscured by dust — does radar alone give safe fallback?)."},
		},
	})

	createSlide(database, s3.ID, "Edge Compute & Data Pipeline Constraints", "CONTENT", 7, map[string]interface{}{
		"type":      "content",
		"heading":   "Edge Compute & Data Pipeline Constraints",
		"introText": "A mine CAS is not a cloud product. Every safety-critical decision runs at the vehicle edge, on constrained hardware, in a harsh physical environment, with no assumption of connectivity. These constraints shape the whole architecture.",
		"items": []map[string]interface{}{
			{"title": "Latency Budget", "description": "End-to-end sensor-to-alert latency must be under 100 ms for safety-critical events. Anything higher and the alert arrives after the operator would already be reacting. Zero tolerance for cloud-round-trip latency — every inference happens on-board."},
			{"title": "Compute Budget", "description": "Rugged embedded systems — Nvidia Jetson Orin / AGX, Qualcomm automotive-grade SoCs, specialised ISO 26262-rated compute modules. Not datacenter GPUs. Power envelope typically 50–200 W per vehicle, passively cooled or fan-cooled with dust filtration."},
			{"title": "Thermal Constraints", "description": "Mining trucks operate in ambient >45°C and direct sun load on cab and bodywork push electronics enclosures past 70°C. Compute modules must be rated for sustained high-temperature operation. Consumer-grade GPUs throttle hard in these conditions and are not viable."},
			{"title": "Connectivity Assumptions", "description": "Mine sites have patchy connectivity: 900 MHz LTE-like radio, leaky-feeder cable in underground, some 5G private networks at newer sites, often Wi-Fi dead zones across the pit. System must be fully operational with zero connectivity — uplink is for logging and updates, never real-time safety."},
			{"title": "Data Pipeline", "description": "Sensor capture → time synchronisation → per-sensor preprocessing → fusion → object detection and tracking → zone evaluation → alert arbitration → operator display/audio and (for L9) intervention command. Entire pipeline must run deterministically in <100 ms, every cycle, for the life of the vehicle."},
			{"title": "Logging", "description": "Every detection, track, alert, and operator response is logged for post-incident analysis. Typical: on-board SSD with 30–90 days rolling buffer, synced to cloud whenever connectivity is available. Logs are evidence — data integrity, time synchronisation, and tamper-evident storage are engineering requirements, not nice-to-haves."},
		},
	})

	// ── Section 4: What the Data Enables ────────────────────────────────────────
	s4 := createSection(database, module.ID, "What the Data Enables", 4)

	createSlide(database, s4.ID, "What the Data Enables", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "What the Data Enables",
		"subtitle": "L8/L9 doesn't just detect threats — it generates the diagnostic data to eliminate them",
		"bullets": []string{
			"Every proximity event maps to a specific L1-7 control gap",
			"AI analysis of event patterns turns raw data into actionable L1-7 recommendations",
			"The system that closes the loop: technology improves the fundamentals it depends on",
			"L1-7 optimisation moves from reactive (post-incident) to proactive (data-driven continuous improvement)",
		},
	})

	createSlide(database, s4.ID, "L8 Event Data as a Diagnostic for L1–7 Gaps", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "L8 Event Data as a Diagnostic for L1–7 Gaps",
		"introText": "The technology layer is not just the last line of defence. It is the most comprehensive data source available on where the L1-7 control hierarchy is failing — if you know how to read it.",
		"items": []map[string]interface{}{
			{
				"title":       "What every L8 proximity event records",
				"description": "Location, time, object types involved, heading geometry, closing speed, and operator response. Across a fleet of 50 trucks operating for 12 months, this is millions of data points — a continuous, granular audit of where vehicle interactions are actually occurring on this site.",
			},
			{
				"title":       "Pattern analysis reveals the L1-7 gaps",
				"description": "Proximity events cluster at a specific intersection → L1 mine design problem (the intersection geometry creates recurring interaction scenarios). Events cluster around a specific operator → L6 training problem. Events spike at shift changeover → L5 procedure problem. Events increase in wet season → L4 vehicle standards or road quality issue. Each pattern points to a specific layer failure.",
			},
			{
				"title":       "From investigation to prevention",
				"description": "Without this data, L1-7 gap analysis relies on incident investigation — which means waiting for something to go wrong. With L8 event data, the gaps are visible before they cause incidents. A dangerous intersection shows up in the data months before a fatality. The intervention can be physical redesign rather than a post-incident inquiry.",
			},
			{
				"title":       "The feedback loop the industry doesn't talk about",
				"description": "Current CAS narrative focuses on detection and intervention — the technology responding to threats. The more powerful argument is that the technology, properly instrumented and analysed, actively improves the L1-7 controls it depends on. A site with good L8 event data analysis will have fewer L8 events over time — not because the system gets better at detecting, but because the hazards it detects get eliminated.",
			},
		},
		"keyFact": "The technology layer pays back more than it costs — because the data it generates improves the entire control hierarchy below it.",
	})

	createSlide(database, s4.ID, "The AI Recommendation Layer", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "The AI Recommendation Layer",
		"introText": "Manual analysis of L8 event data by a safety manager reviewing weekly charts works — but it hits cognitive limits quickly. Pattern recognition across thousands of events, machine types, operators, locations, and time periods exceeds human capacity at fleet scale. AI agents change what is possible.",
		"items": []map[string]interface{}{
			{
				"title":       "From dashboards to recommendations",
				"description": "The shift from charts that safety managers interpret to AI agents that generate specific recommendations is not incremental — it changes the nature of the output. Instead of 'here is the data', the system says 'here is what to do about it, at this location, for this operator, against this specific L1-7 layer'.",
			},
			{
				"title":       "What AI-generated L1-7 recommendations look like",
				"description": "Examples: 'The intersection at grid X3-Y7 accounts for 18% of total proximity events despite representing 3% of haul road usage — a traffic island or one-way section would eliminate the interaction geometry' (L1/L3 recommendation). 'Operator #4421 generates 4.2x the fleet-average Zone 2 events at haul road curves — recommend targeted training review' (L6). 'Proximity events at the crusher approach spike 40% in the hour before blast — recommend traffic management plan review for blast window' (L5).",
			},
			{
				"title":       "Continuous improvement without additional inspections",
				"description": "AI-driven L1-7 analysis runs continuously across the full event history. Every new shift's data refines the model. Recommendations are updated as patterns evolve — new pit areas generate new patterns, and the system adapts. This is proactive safety management at a scale no human inspection programme can match.",
			},
		},
	})

	createSlide(database, s4.ID, "The Continuous Improvement Loop", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Continuous Improvement Loop",
		"description": "The perception-based CAS is not a static safety layer — it is the data source for a continuous improvement cycle that progressively reduces the hazards it was deployed to manage.",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"title": "L8/L9 detects interactions", "description": "Every proximity event, advisory alert, and intervention is timestamped and logged with full context — location, object types, headings, speeds, operator response."},
				{"title": "AI agents analyse event patterns", "description": "Pattern recognition across the full event history identifies structural contributors: specific locations, operators, times of day, weather conditions, or operational sequences that repeatedly generate interactions."},
				{"title": "Specific L1-7 recommendations generated", "description": "Each identified pattern is mapped to the relevant EMESRT level — mine design gap (L1-3), procedure gap (L4-6) — and a concrete recommendation is generated for that specific gap."},
				{"title": "L1-7 controls improved", "description": "The site implements the recommended control change: physical redesign, procedural update, targeted training, or supervisor intervention. The change is logged against the pattern that triggered it."},
				{"title": "L8/L9 event volume decreases", "description": "As L1-7 controls improve, the interactions that generated L8 events are eliminated or reduced. The system measures its own effectiveness. Fewer L8 events means the loop is working — and builds the case for the technology investment."},
			},
		},
	})

	// ── Section 5: Product Requirements ────────────────────────────────────────
	s5 := createSection(database, module.ID, "Product Requirements", 5)

	createSlide(database, s5.ID, "Product Requirements — What the System Must Do", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Product Requirements — What the System Must Do",
		"subtitle": "High-level requirements for a perception-based CAS meeting NSW TRG Layer 8/9 and EMESRT L7/8/9",
		"bullets": []string{
			"Functional requirements: detect, classify, track, and intervene",
			"Performance requirements: probability of detection, false positive rate, latency, availability",
			"Development priority order: build incrementally, validate at each stage before proceeding",
		},
	})

	createSlide(database, s5.ID, "Functional Requirements", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Functional Requirements",
		"introText": "These nine functional requirements define what a perception-based CAS must do to comply with NSW TRG Layer 8/9 and EMESRT Levels 7, 8, and 9.",
		"items": []map[string]interface{}{
			{
				"title":       "FR-1 — Tagless detection",
				"description": "Detect and classify any physical object (vehicle, person, static obstruction) within Zone 2 (15–30m) without requiring the object to carry any electronic tag. This is the defining requirement that differentiates perception-based from transponder-based systems.",
			},
			{
				"title":       "FR-2 — Trajectory tracking and collision prediction",
				"description": "Track object trajectory and predict collision probability within Zone 2 to generate advisory alerts calibrated to closing speed and heading geometry. The alert must convey object type and threat vector — not just proximity.",
			},
			{
				"title":       "FR-3 — Machine intervention",
				"description": "Trigger speed reduction and controlled braking for Zone 1 entries (<15m) when the operator has not responded to the Zone 2 advisory within the reaction time window. Intervention must be a controlled deceleration profile, not an emergency stop.",
			},
			{
				"title":       "FR-4 — Environmental robustness",
				"description": "Maintain detection capability in mining environmental conditions: dust density up to 50mg/m³, temperature -20°C to +70°C operating, vibration per ISO 16750-4 Class 4, IP67 housing, and continuous 24-hour operation.",
			},
			{
				"title":       "FR-5 — EMESRT L7 HMI",
				"description": "Operator display showing zone status, object type, and directional threat indicator — designed for the cognitive load environment of a haul truck cab. Display layout must comply with EMESRT Level 7 HMI storyboard requirements.",
			},
			{
				"title":       "FR-6 — System health self-monitoring",
				"description": "Continuous self-diagnostic with automatic degraded-mode alert when any sensor falls below performance threshold. Silent failure — where the system stops working without alerting — is not an acceptable failure mode.",
			},
			{
				"title":       "FR-7 — Event data recording",
				"description": "Timestamped logging of all detection events, alert triggers, intervention activations, and sensor health transitions — retained for a minimum of 30 days and retrievable for post-incident review and L1-7 gap analysis.",
			},
			{
				"title":       "FR-8 — Operator override and acknowledgement logging",
				"description": "Every operator acknowledgement of a Zone 2 advisory and every Zone 1 intervention override must be logged with timestamp, vehicle ID, and operator ID. This audit trail is required for safety case compliance and for identification of repeat override patterns.",
			},
			{
				"title":       "FR-9 — Fleet management system integration",
				"description": "Real-time system health telemetry to the site FMS for compliance dashboard visibility. Shift-level performance reports (detection count, NVA rate, availability) exportable in formats compatible with NSW TRG audit requirements.",
			},
		},
	})

	createSlide(database, s5.ID, "Performance Requirements", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Performance Requirements",
		"introText": "These six performance requirements define the quantitative thresholds the system must meet in production operations. All are measurable and must be validated against production data, not test-scenario results alone.",
		"items": []map[string]interface{}{
			{
				"title":       "PR-1 — Detection probability",
				"description": ">99.5% for person in Zone 2 (15–30m) in normal conditions; >95% in adverse conditions (dust, rain). Measured per EMESRT PR-5A test scenarios and validated in production over >500 machine-hours.",
			},
			{
				"title":       "PR-2 — False positive rate",
				"description": "<2 NVA alerts per machine-hour for advisory (Zone 2); <0.1 per machine-hour for intervention triggers (Zone 1). Measured in production over >500 machine-hours. Exceeding these thresholds triggers system review before continued operation.",
			},
			{
				"title":       "PR-3 — Latency",
				"description": "Detection-to-alert <200ms; detection-to-intervention initiation <500ms. Measured end-to-end from sensor ingestion to vehicle control signal. Latency budget must be allocated across sensor frame rate, fusion processing, CAN bus transmission, and brake actuator response.",
			},
			{
				"title":       "PR-4 — System availability",
				"description": ">99.5% during scheduled operating time. Includes sensor cleaning cycle downtime budget. Availability monitoring must be continuous and reported per shift.",
			},
			{
				"title":       "PR-5 — Classification accuracy",
				"description": ">97% correct object type classification (vehicle/person/static) measured against PR-5A scenario storyboards under production conditions. Classification errors must be logged and reviewed — they are the primary input to model improvement cycles.",
			},
			{
				"title":       "PR-6 — Intervention behaviour",
				"description": "Controlled deceleration profile calibrated to detection certainty level — not an emergency stop. A false intervention triggered at high certainty threshold must itself not cause loss of vehicle control. Deceleration profiles must be validated with OEM braking specifications.",
			},
		},
	})

	createSlide(database, s5.ID, "Development Priority Order", "PROCESS", 3, map[string]interface{}{
		"type":        "process",
		"heading":     "Development Priority Order",
		"description": "A perception-based CAS must be built incrementally. Each stage validates the foundation before the next layer is added. Attempting classification before detection is reliable, or intervention before alarm fatigue is solved, produces unsafe outcomes.",
		"steps": []map[string]interface{}{
			{
				"title":       "Step 1 — Sensor platform and data collection",
				"tier":        "Operate",
				"description": "Mount the multi-modal sensor stack on the target vehicle class. Collect >1,000 hours of labelled operating data across target conditions. Establish sensor fusion baseline. Do not attempt classification yet — this step is about instrumentation quality.",
			},
			{
				"title":       "Step 2 — Detection and tracking",
				"tier":        "Operate",
				"description": "Achieve reliable object detection and tracking in Zone 2 conditions without classification. Establish the false positive baseline before adding ML classification — understand what the raw sensor fusion layer produces.",
			},
			{
				"title":       "Step 3 — Classification",
				"tier":        "React",
				"description": "Train and validate person vs vehicle vs static object classification with >97% accuracy. Validate against PR-5A scenario storyboards in controlled conditions. Conduct adversarial testing: low-visibility, dust, night, glare, and partial occlusion scenarios.",
			},
			{
				"title":       "Step 4 — Zone integration and alert logic",
				"tier":        "React",
				"description": "Integrate zone-based alert triggering with dynamic zone sizing by closing speed. Implement EMESRT L7/L8 compliant HMI integration. Conduct alarm fatigue testing: measure NVA rate in production operations and tune until <2 NVA alerts per machine-hour. This is the gate before Step 5.",
			},
			{
				"title":       "Step 5 — Machine intervention",
				"tier":        "React",
				"description": "Integrate with vehicle CAN bus for EMESRT Level 9 intervention. Implement intervention logic and failsafe validation. Complete OEM integration or aftermarket ECU regulatory approval. Pursue NSW TRG Layer 9 compliance pathway.",
				"warnings": []string{
					"Do not attempt machine intervention until Steps 1–4 are validated in production operations across >500 machine-hours. An unsafe intervention — triggered on a non-genuine threat — is more dangerous than no system.",
				},
			},
		},
	})

	// ── Section 6: Why This Enables What Wasn't Possible Before ────────────────
	s6 := createSection(database, module.ID, "Why This Enables What Wasn't Possible Before", 6)

	createSlide(database, s6.ID, "What Perception Makes Possible", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "What Perception Makes Possible",
		"subtitle": "The specific capability gaps that perception-based architecture closes",
		"bullets": []string{
			"Untagged hazard detection — the blind spot that transponder systems cannot address",
			"Alarm fatigue elimination — contextual classification vs zone-entry triggering",
			"Head-to-head comparison: perception vs transponder across operational dimensions",
		},
	})

	createSlide(database, s6.ID, "The Untagged Hazard Problem — Solved", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Untagged Hazard Problem — Solved",
		"introText": "Perception does not ask whether an object carries a tag. It asks whether an object is physically present. This changes the risk profile for every scenario where tagging compliance is imperfect — which is every site, every shift.",
		"items": []map[string]interface{}{
			{
				"title":       "Contractor arriving on site without tag configured",
				"description": "Happens on every large site, every day. The visitor vehicle management process fails silently — the vehicle enters the pit without a functioning tag, and a transponder CAS has zero knowledge of its presence. A perception-based system detects it by the physical fact of its being there.",
			},
			{
				"title":       "Maintenance crew working in pit without PPE tag online",
				"description": "A tradesperson working on a conveyor transfer point in the pit, whose hard hat tag battery died two shifts ago. For a transponder system, this person is invisible to every vehicle operating in that zone. For a perception system, the person is a detected and classified hazard.",
			},
			{
				"title":       "Rock and spoil spills, mis-parked equipment, third-party vehicles",
				"description": "A spoil spill partially blocking the haul road carries no tag. A piece of equipment parked in a non-designated area does not have an active transponder signal. A survey vehicle from an adjacent property has no site tag configured. All are visible to perception; all are invisible to transponder.",
			},
			{
				"title":       "The residual risk calculation",
				"description": "The transponder system's risk model assumes zero untagged hazards. The real-world risk model on any large site has a non-zero untagged hazard rate on every shift. Perception eliminates this entire hazard category from the residual risk calculation.",
			},
		},
	})

	createSlide(database, s6.ID, "The Alarm Fatigue Problem — Structurally Addressed", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "The Alarm Fatigue Problem — Structurally Addressed",
		"introText": "Perception-based ML classification does not merely detect proximity — it classifies threat context. This is the architectural change that produces fundamentally lower NVA alarm rates. Combined with operational memory as the system matures, NVA rates trend toward zero over a deployment lifetime.",
		"items": []map[string]interface{}{
			{
				"title":       "The classification difference",
				"description": "A transponder CAS generates an alert when a tag enters the zone radius — it cannot distinguish a vehicle on a parallel road (no threat) from a vehicle on a collision heading (genuine threat). Both are in the zone; both generate the same alert. A perception-based system classifies not just proximity but heading, velocity vector, and trajectory — only genuine threats alarm.",
			},
			{
				"title":       "Background objects filtered at the classification layer",
				"description": "Pit walls, berms, passing vehicles on separated haul roads, stationary equipment in designated parking, and slow-moving support vehicles in non-collision trajectories are classified as non-threats before the alert decision is made. They appear in Zone 3 (awareness only) and do not escalate to Zone 2 advisory.",
			},
			{
				"title":       "The NVA rate improvement trajectory",
				"description": "Current automotive-grade ML classification models achieve NVA rates below 1 per machine-hour on mining-environment data — a 50x+ reduction from the 50,000-alarms-on-two-dozers benchmark. As operational memory develops and the system learns site-specific patterns (switchback locations, normal queue behaviour, bench geometry), NVA rates continue to decrease.",
			},
			{
				"title":       "The trust restoration pathway",
				"description": "Operators who experienced Gen 1 alarm fatigue can be shown production NVA rate data trending toward zero. The metric is concrete, measurable, and visible. A system that demonstrably reduces NVA rates in its first 30 days of operation earns operator trust by evidence — not by instruction.",
			},
		},
	})

	createSlide(database, s6.ID, "Perception-Based vs Transponder-Based", "COMPARISON", 3, map[string]interface{}{
		"type":        "comparison",
		"heading":     "Perception-Based vs Transponder-Based CAS",
		"description": "A fair architectural comparison. Both approaches can achieve EMESRT Level 9 — the question is which hazard classes each can detect, and at what operational cost.",
		"items": []map[string]interface{}{
			{
				"name":        "Perception-Based CAS",
				"description": "Camera + 4D radar + LiDAR sensor fusion. Detects any physical object regardless of electronic tagging.",
				"pros": []string{
					"Detects any physical object — tagged or untagged, known or unknown to the system",
					"Classifies object type (vehicle / person / static) via on-device ML inference",
					"4D imaging radar operates in dust, rain, fog, and GNSS-denied pit environments",
					"ML-based contextual classification targets <2 NVA advisory alerts per machine-hour",
					"Zero tag compliance overhead — sensors observe the physical world directly",
					"Naturally immune to multi-bench false positives — cannot see through pit walls",
				},
				"cons": []string{
					"Higher hardware cost than transponder systems — automotive-grade edge SoCs (NVIDIA Orin)",
					"Requires mining-specific training data programme for full classification performance",
					"Camera and LiDAR performance degrades in heavy mineral dust — radar-primary mode required",
					"More complex system integration — safety case and OEM approval path required for Level 9",
					"Less proven at scale than transponder systems — still in production validation stage",
				},
			},
			{
				"name":        "Transponder-Based CAS",
				"description": "Tag-based proximity detection. Every hazard must carry an active electronic transponder. Proven at scale globally.",
				"pros": []string{
					"Proven at scale — 15,000+ vehicles globally (includes Level 7 through Level 9 variants)",
					"Mature OEM and FMS integration ecosystems",
					"UWB-based implementations provide sub-10cm positional accuracy for confined areas",
					"Pedestrian-worn device gives the person an independent alert — defence-in-depth",
					"Lower unit cost per installation; simpler safety case for advisory-only (Level 7/8) deployment",
					"Well-understood regulatory acceptance pathway in most jurisdictions",
				},
				"cons": []string{
					"Only detects hazards carrying an active, functioning electronic tag — untagged objects are completely invisible",
					"Requires 100% tag compliance across all site personnel and vehicles, every shift",
					"Tag failure modes (dead battery, forgotten PPE, interference) create invisible hazards",
					"Structural false positives from multi-bench operations — no model of 3D mine geometry",
					"Structural false positives from normal mine manoeuvres (switchbacks, tight turns) — no mine geometry context",
					"Tag maintenance, replacement, and compliance auditing is a sustained operational overhead",
				},
			},
		},
	})

	// ── Section 7: Implementation Path ─────────────────────────────────────────
	s7 := createSection(database, module.ID, "Implementation Path", 7)

	createSlide(database, s7.ID, "From Lab to Pit", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "From Lab to Pit — The Implementation Path",
		"subtitle": "Phased deployment from data collection to fleet-wide Level 9",
		"bullets": []string{
			"Five overlapping phases spanning 18–24 months to full Level 9 deployment",
			"Each phase gates the next — NSW TRG compliance pathway built in",
			"Change management is the most underestimated requirement",
		},
	})

	createSlide(database, s7.ID, "From Lab to Pit — Deployment Phases", "PROCESS", 1, map[string]interface{}{
		"type":        "process",
		"heading":     "From Lab to Pit — Deployment Phases",
		"description": "Perception-based CAS deployment follows five overlapping phases. The key constraint is the validation gate between each phase.",
		"steps": []map[string]interface{}{
			{
				"title":    "Phase 1 — Data collection and sensor validation",
				"duration": "Months 1–6",
				"tier":     "Operate",
				"description": "Multi-modal sensor rig mounted on target vehicle class. Collect and label ground-truth data across all target operating conditions. Establish sensor performance baseline against PR-5A scenarios.",
			},
			{
				"title":    "Phase 2 — Detection and classification model development",
				"duration": "Months 4–12",
				"tier":     "Operate",
				"description": "Train and validate object detection models on mining-specific datasets, supplemented by synthetic data and automotive transfer learning. Achieve >97% classification accuracy in controlled test conditions. Begin site trials in alert-only mode.",
			},
			{
				"title":    "Phase 3 — Alert system integration and alarm fatigue testing",
				"duration": "Months 10–18",
				"tier":     "React",
				"description": "Integrate with onboard display per EMESRT L7/L8 HMI storyboards. Measure NVA alarm rate in production operations across a representative vehicle fleet. Tune classification thresholds and zone sizing until NVA rate is below 2 per machine-hour. This is the gate for NSW TRG Layer 8 compliance pathway.",
			},
			{
				"title":    "Phase 4 — Machine intervention development, validation, and regulatory approval",
				"duration": "Months 16–24",
				"tier":     "React",
				"description": "CAN bus integration with vehicle drivetrain for EMESRT Level 9. Intervention logic, failsafe validation, and adversarial testing. OEM acceptance letter required by most Tier 1 operators before intervention systems may be fitted. Safety case per ISO 13849 or IEC 62061. Phase 3 NVA validation is a hard prerequisite — do not begin Phase 4 until Phase 3 is complete.",
			},
			{
				"title":    "Phase 5 — Site-wide deployment and continuous improvement",
				"duration": "Months 18+",
				"tier":     "React",
				"description": "Fleet rollout across all target vehicle classes. L8 event data analysis programme commences — AI-driven L1-7 gap identification and recommendation generation. Model updates as operating condition data accumulates. Performance reporting infrastructure to support NSW TRG audit requirements.",
			},
		},
	})

	createSlide(database, s7.ID, "NSW TRG Compliance — What the Regulator Will Audit", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "NSW TRG Compliance — What the Regulator Will Audit",
		"introText": "The NSW Resources Regulator's enforcement posture on vehicle interaction is active and unannounced. Sites deploying CAS technology must be prepared to demonstrate compliance evidence across all nine layers — not just the technology layers.",
		"items": []map[string]interface{}{
			{
				"title":       "Evidence that Layers 1–7 are in place",
				"description": "This is not optional. The Regulator will inspect physical segregation infrastructure, check traffic management plan currency, review operator training records, and verify supervisor competency. Technology deployed without Layers 1–7 does not satisfy the 9-layer requirement and will attract enforcement action.",
			},
			{
				"title":       "System performance data",
				"description": "Detection probability per hazard class, NVA alarm rate per machine-hour, system availability per shift, and intervention event logs — all logged, timestamped, and auditable. The Regulator expects to see this data at inspection.",
			},
			{
				"title":       "Change management evidence",
				"description": "Operator training records for the CAS system, supervisor sign-off on HMI competency, near-miss and incident data before and after deployment, and evidence that operators have been involved in system calibration.",
			},
			{
				"title":       "Vendor compliance documentation",
				"description": "PR-5A compliance against specific EMESRT scenario storyboards, supported by test data from an accredited testing facility. For Level 9 intervention systems: OEM acceptance documentation, safety case per ISO 13849 or IEC 62061, and records of jurisdictional inspectorate acceptance.",
			},
		},
	})

	createSlide(database, s7.ID, "Change Management — The Most Underestimated Requirement", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Change Management — The Most Underestimated Requirement",
		"introText": "Technical performance is necessary but not sufficient. A CAS that operators do not trust, understand, or engage with will not prevent fatalities — regardless of its sensor specifications. Change management is the make-or-break variable in CAS deployments.",
		"items": []map[string]interface{}{
			{
				"title":       "The Gen 1 legacy barrier",
				"description": "Operators who experienced alarm fatigue in Gen 1 deployments have a rational prior belief that CAS systems are noise generators. Overcoming this requires demonstrated evidence, not instructions. The first 30 days of a new deployment are the credibility window — if NVA rate is not visibly low in that window, operator trust does not recover.",
			},
			{
				"title":       "Operator involvement in calibration",
				"description": "Giving operators a mechanism to report NVA alarms in the first 30 days creates ownership and provides the highest-quality NVA feedback data available. Operators know better than any algorithm what constitutes a non-genuine alert in their specific operating environment. Involvement generates trust; trust generates engagement.",
			},
			{
				"title":       "Shift supervisor training and visible metrics",
				"description": "Supervisors must understand what advisory alerts mean, what an intervention event log shows, and how to read the system health dashboard. Visible performance metrics posted in the crib room — NVA rate trending down, detection probability holding above target — provide proof the system is working and improving.",
			},
			{
				"title":       "The Level 9 activation gate",
				"description": "Never activate machine intervention before operators have experienced Level 7/8 advisory performance for a minimum of three months. Operators who understand and trust the advisory system accept intervention as a logical extension. Operators who do not trust the advisory system treat intervention as a hostile, unpredictable machine behaviour — and begin circumventing it.",
			},
		},
	})

	// ── Section 8: Knowledge Check ──────────────────────────────────────────────
	s8 := createSection(database, module.ID, "Knowledge Check", 8)

	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Collision Avoidance Systems — Perception-Based Approach Assessment",
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
			question:    "What percentage of mining industry fatalities are attributable to vehicle interaction failures each year?",
			explanation: "ICMM Health & Safety Performance data consistently places vehicle interaction failures at approximately 25–35% of annual mining fatalities. This figure has persisted despite decades of CAS investment — indicating a structural problem with existing technology architectures, not merely insufficient investment.",
			options: []option{
				{"10–15%", false},
				{"20–25%", false},
				{"25–35%", true},
				{"50–60%", false},
			},
		},
		{
			question:    "In the EMESRT 9-level framework, which levels fall under the Design tier?",
			explanation: "Levels 1–3 are Design controls — built into the physical and procedural fabric of the mine before operations begin. Level 1: Site Requirements. Level 2: Segregation Controls. Level 3: Physical & Procedural Design Controls. Levels 4–6 are Operate controls. Levels 7–9 are React (technology) controls.",
			options: []option{
				{"Levels 1 and 2 only", false},
				{"Levels 1, 2, and 3", true},
				{"Levels 1 through 6", false},
				{"All levels contribute to design", false},
			},
		},
		{
			question:    "In the NSW TRG ROVOA 9-layer defence model, which layers contain Collision Avoidance Technology (CAT)?",
			explanation: "The NSW TRG ROVOA explicitly places CAT — proximity detection and machine intervention — in Layers 8 and 9. Layers 1–7 are fundamental design and operational controls that must be in place before technology layers are effective.",
			options: []option{
				{"Layers 1 and 2", false},
				{"Layers 5, 6, and 7", false},
				{"Layers 8 and 9", true},
				{"All 9 layers contribute to CAT effectiveness", false},
			},
		},
		{
			question:    "What is the fundamental limitation of transponder-based collision avoidance systems that perception-based systems solve?",
			explanation: "Any object not carrying an active, functioning tag is completely invisible to a transponder-based CAS. This includes contractors without site tags, failed tags, maintenance equipment in transit, rock falls, and debris. Perception-based systems detect any physical object regardless of tagging.",
			options: []option{
				{"Transponder systems cannot work underground", false},
				{"Transponder systems require GNSS connectivity at all times", false},
				{"Transponder systems only detect hazards that carry an active electronic tag — untagged objects are invisible", true},
				{"Transponder systems cannot trigger Level 9 machine intervention", false},
			},
		},
		{
			question:    "In the Glencore/Wabtec Gen 1 deployment, 50,000 alerts were recorded on two dozers. What operational consequence did this cause?",
			explanation: "When a system generates overwhelming volumes of non-value-add alarms, operators rationally discount all alerts — including genuine threats. A system that is wrong 49,999 times makes the 50,000th alert invisible. Alarm fatigue is identified as the primary CAS deployment failure mode.",
			options: []option{
				{"The dozers were taken out of service", false},
				{"Alarm fatigue caused operators to begin ignoring the system — including genuine threat alerts", true},
				{"The system was upgraded to Level 9 to automatically handle the alerts", false},
				{"Glencore sued Wabtec for product non-compliance", false},
			},
		},
		{
			question:    "A perception-based CAS consistently generates Zone 2 alerts at a hairpin ramp on the descent from the upper pit. No collision events have occurred at this location and sensor diagnostics show no faults. The most likely root cause is:",
			explanation: "At the apex of a switchback or hairpin turn, a haul truck's kinematic velocity vector points directly at the oncoming lane. A system projecting linear paths from current velocity flags this as a collision course every single time. The truck is executing a normal, road-geometry-constrained manoeuvre — but the system has no model of mine road geometry to distinguish this from an unexpected heading change. This is a structural false positive that fires on every turn, every shift.",
			options: []option{
				{"Dust is causing false sensor returns at this location", false},
				{"The system is projecting linear paths from current kinematics and lacks mine geometry context to distinguish a normal ramp turn from an unexpected heading change", true},
				{"The haul road at the switchback is outside the configured detection zone", false},
				{"Zone 2 alert thresholds are misconfigured for the site", false},
			},
		},
		{
			question:    "A dozer working on bench floor 1 and a haul truck on bench 3 directly above are generating repeated false proximity alerts on a GNSS-based CAS. The most accurate explanation is:",
			explanation: "GNSS-based positioning systems evaluate collision risk in 2D or with poor vertical resolution. The dozer and truck appear co-located in 2D coordinates — but a pit wall and 20-30 metres of vertical separation mean they cannot physically interact. The positioning system has no model of pit wall or bench geometry. A perception-based system cannot see through pit walls — the physical mine geometry filters this false positive automatically.",
			options: []option{
				{"A GNSS signal error is placing the truck at the wrong location", false},
				{"The GNSS system evaluates proximity in 2D and has no model of bench geometry or pit walls — machines physically separated by vertical distance appear co-located", true},
				{"The dozer's transponder tag is broadcasting at excessive range", false},
				{"The detection zone radius should be reduced for bench operations", false},
			},
		},
		{
			question:    "Which sensor modality maintains reliable detection simultaneously in heavy dust, rain, and GNSS-denied pit environments?",
			explanation: "4D imaging radar uses FMCW radio waves, which are not materially attenuated by dust or precipitation, and do not depend on satellite signals. LiDAR and cameras degrade significantly in heavy mineral dust. RTK GNSS requires satellite reception. Only radar addresses all three conditions simultaneously.",
			options: []option{
				{"AI vision cameras — widest field of view", false},
				{"Solid-state LiDAR — highest range resolution", false},
				{"4D imaging radar — radio waves unaffected by particulates and satellite denial", true},
				{"RTK GNSS — highest position accuracy", false},
			},
		},
		{
			question:    "What is 'operational memory' in the context of CAS, and why does its absence cause persistent false positives?",
			explanation: "Operational memory is a system's accumulated model of what normal looks like on a specific mine site — where machines normally travel, when they turn, what approach patterns precede a dump. Without it, every event is evaluated in isolation. A switchback turn generates the same alert on day 500 as on day 1 because the system has no knowledge that trucks always turn here. False positive rates remain structurally fixed rather than improving over time.",
			options: []option{
				{"The system's audit log — required for post-incident review under the NSW TRG", false},
				{"A system's accumulated model of normal operational patterns on this mine — its absence means every event is evaluated in isolation and false positive rates never improve", true},
				{"The FMS dispatch history used to predict truck destinations", false},
				{"The operator's memory of previous alert events — relevant to change management", false},
			},
		},
		{
			question:    "L8/L9 event data that shows 28% of all proximity alerts cluster at a single haul road intersection during shift changeover most directly indicates:",
			explanation: "A cluster of proximity events at a specific location during shift changeover points to two distinct gap types: a Layer 5 (operating procedures) issue — the shift changeover traffic management protocol is generating mixed-traffic interactions — and potentially a Layer 1/3 issue if the intersection geometry itself creates the conflict. The data is a diagnostic pointing to specific L1-7 controls to investigate, before an incident occurs.",
			options: []option{
				{"A sensor calibration issue specific to that location", false},
				{"A Layer 5 (operating procedures) and potentially Layer 1/3 (design) gap — the intersection generates recurring interactions that a procedure change or physical redesign could eliminate", true},
				{"The NVA rate is too high and the zone threshold needs tuning", false},
				{"An alarm fatigue problem — operators are not responding to alerts at this location", false},
			},
		},
		{
			question:    "A mine's transponder-based CAS passes its PR-5A compliance test with no failures. Three months later, a pedestrian working in the pit without a functioning PPE tag is struck by a haul truck that received no alert. Which statement best explains this outcome?",
			explanation: "PR-5A validates detection performance against tagged hazards in controlled test scenarios. An untagged pedestrian is architecturally invisible to a transponder system regardless of PR-5A compliance status. No sensor failure, no operator error, and no Layers 1–7 gap explains this — it is the inherent limitation of tag-dependent architecture.",
			options: []option{
				{"A sensor hardware fault — PR-5A compliance cannot prevent equipment failures", false},
				{"PR-5A compliance tests tagged hazards in controlled scenarios — it cannot detect untagged objects, which is the architecture's inherent limitation", true},
				{"Alarm fatigue — the operator had received too many alerts and ignored this one", false},
				{"The mine failed to implement Layers 1–7 before deploying technology", false},
			},
		},
		{
			question:    "A site prepares for an unannounced NSW Regulator inspection after deploying a CAT Layer 8 proximity alert system. The site has no formal traffic management plan and mixed light/heavy vehicle haul roads. What is the most likely inspection outcome?",
			explanation: "The NSW Regulator's enforcement position is explicit: technology deployed without Layers 1–7 in place does not constitute compliance. Missing traffic management plan (Layer 5) and mixed-traffic roads (Layer 3 failure) are fundamental control gaps. Layer 8 technology investment does not offset Layer 1–7 deficiencies — prohibition notices have been issued in this scenario.",
			options: []option{
				{"The Layer 8 technology will be credited as demonstrating safety commitment", false},
				{"The site will receive a minor improvement notice with a 90-day rectification period", false},
				{"A prohibition notice is likely — technology without Layers 1–7 does not satisfy the regulatory requirement", true},
				{"The TRG applies to new mines only — existing operations have a compliance transition period", false},
			},
		},
	}

	for i, q := range questions {
		qModel := models.QuizQuestion{
			ID:        uuid.New().String(),
			QuizID:    quiz.ID,
			Type:      models.QuestionMultipleChoice,
			Question:  q.question,
			Points:    1,
			SortOrder: i,
		}
		if q.explanation != "" {
			qModel.Explanation = &q.explanation
		}
		database.Create(&qModel)

		for j, o := range q.options {
			database.Create(&models.QuizOption{
				ID:         uuid.New().String(),
				QuestionID: qModel.ID,
				Text:       o.text,
				IsCorrect:  o.isCorrect,
				SortOrder:  j,
			})
		}
	}

	createSlide(database, s8.ID, "CAS Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Collision Avoidance Systems — Perception-Based Approach Assessment",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{
				"question":     "What percentage of mining industry fatalities are attributable to vehicle interaction failures each year?",
				"options":      []string{"10–15%", "20–25%", "25–35%", "50–60%"},
				"correctIndex": 2,
				"explanation":  "ICMM Health & Safety Performance data consistently places vehicle interaction failures at approximately 25–35% of annual mining fatalities.",
			},
			{
				"question":     "In the EMESRT 9-level framework, which levels fall under the Design tier?",
				"options":      []string{"Levels 1 and 2 only", "Levels 1, 2, and 3", "Levels 1 through 6", "All levels"},
				"correctIndex": 1,
				"explanation":  "Levels 1–3 are Design controls. Level 3 (Physical & Procedural Design Controls) is Design, not Operate. Levels 4–6 are Operate. Levels 7–9 are React.",
			},
			{
				"question":     "In the NSW TRG ROVOA 9-layer defence model, which layers contain CAT (Collision Avoidance Technology)?",
				"options":      []string{"Layers 1 and 2", "Layers 5, 6, and 7", "Layers 8 and 9", "All 9 layers"},
				"correctIndex": 2,
				"explanation":  "CAT — proximity detection and machine intervention — sits in Layers 8 and 9. Layers 1–7 are fundamental design and operational controls.",
			},
			{
				"question":     "What is the fundamental limitation of transponder-based CAS that perception-based systems solve?",
				"options":      []string{"Cannot work underground", "Requires GNSS at all times", "Only detects tagged hazards — untagged objects are invisible", "Cannot trigger Level 9 intervention"},
				"correctIndex": 2,
				"explanation":  "Any untagged object is completely invisible to a transponder CAS. Perception detects any physical object regardless of whether it carries a tag.",
			},
			{
				"question":     "A CAS consistently false-alarms at a hairpin ramp turn. No sensor faults detected. The most likely cause is:",
				"options":      []string{"Dust attenuation at this location", "Linear path projection flagging a normal geometry-constrained turn as a collision course — no mine road model", "Zone radius misconfiguration", "Operator has disabled alerts at this location"},
				"correctIndex": 1,
				"explanation":  "The switchback false positive is structural: the system projects linear paths from current kinematics and has no model of mine road geometry to know the truck is executing a normal turn.",
			},
			{
				"question":     "What is 'operational memory' and why does its absence cause persistent false positives?",
				"options":      []string{"The audit log of past alerts", "An accumulated model of normal mine operational patterns — without it, every event is evaluated in isolation and false positive rates never improve", "FMS dispatch history", "Operator familiarity with the system"},
				"correctIndex": 1,
				"explanation":  "Without operational memory, a switchback fires the same alert on day 500 as day 1. The system has no knowledge of where trucks normally go on this mine.",
			},
			{
				"question":     "L8 event data showing 28% of proximity alerts cluster at one intersection during shift changeover most directly indicates:",
				"options":      []string{"A sensor calibration issue", "A Layer 5 procedure and potential Layer 1/3 design gap — specific L1-7 controls to investigate before an incident occurs", "NVA rate threshold needs tuning", "Alarm fatigue at this location"},
				"correctIndex": 1,
				"explanation":  "Clustered event data is a diagnostic of specific L1-7 control gaps — in this case, shift changeover procedures (L5) and potentially intersection geometry (L1/L3).",
			},
			{
				"question":     "Which sensor modality maintains reliable detection simultaneously in heavy dust, rain, and GNSS-denied pit environments?",
				"options":      []string{"AI vision cameras", "Solid-state LiDAR", "4D imaging radar", "RTK GNSS"},
				"correctIndex": 2,
				"explanation":  "4D imaging radar uses FMCW radio waves, unaffected by dust or precipitation and independent of satellite signals.",
			},
		},
	})

	createSlide(database, s8.ID, "Module Complete", "COMPLETION", 1, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete",
		"message": "You have completed Collision Avoidance Systems — The Perception-Based Approach. You now have a rigorous technical and regulatory foundation for building, deploying, or evaluating a perception-based CAS that meets NSW TRG Layer 8/9 and EMESRT Level 9 requirements.",
		"keyTakeaways": []string{
			"Approximately 25–35% of annual mining fatalities are caused by vehicle interaction failures (ICMM data) — the problem is structural, not incremental",
			"The EMESRT framework has 9 levels: Design (L1–3), Operate (L4–6), React/Technology (L7–9) — Design controls must precede technology",
			"NSW TRG ROVOA codifies the same framework into regulatory requirements — non-compliance attracts prohibition notices, not improvement notices",
			"Transponder systems are blind to untagged hazards — perception detects any physical object regardless of tagging",
			"False positives have three structural sources: proximity without path context, path projection without mine geometry (the switchback problem), and positioning without mine geometry (the multi-bench problem)",
			"Operational memory — a system's accumulated model of normal mine operations — is the missing control that will reduce false positive rates from structurally fixed to continuously improving",
			"L8/L9 event data is the most comprehensive diagnostic of L1-7 control gaps — AI agents can turn this data into specific L1-7 improvement recommendations",
			"4D imaging radar provides all-weather, GNSS-denied detection that cameras and LiDAR cannot match in adverse conditions",
			"Development must proceed in order: detect → classify → validate NVA rate → then and only then activate Level 9 intervention",
		},
		"showCertificate": true,
	})

	log.Println("CAS module seeded with 9 sections and knowledge check")
}
