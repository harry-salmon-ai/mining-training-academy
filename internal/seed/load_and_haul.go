package seed

import (
	"encoding/json"
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedLoadAndHaulModule(database *gorm.DB, force bool) {
	var moduleCount int64
	database.Model(&models.Module{}).Where("slug = ?", "load-and-haul-operations").Count(&moduleCount)
	if moduleCount > 0 {
		if !force {
			log.Println("Load & Haul module already exists, skipping")
			return
		}
		forceDeleteModule(database, "load-and-haul-operations")
	}

	log.Println("Seeding Load & Haul module...")

	// Find category and admin user
	var category models.Category
	database.Where("slug = ?", "mining-operations").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed Load & Haul: category or admin not found")
		return
	}

	desc := "Understanding the engine room of the mining value chain — loading areas, equipment types, truck presentation, hauling, and all dump types including highwalls, paddocks, crusher and road maintenance dumping."
	dur := 90
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Load & Haul Operations",
		Slug:        "load-and-haul-operations",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelFoundation,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"load & haul", "trucks", "excavators", "dumping", "crusher", "haul roads", "safety"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// Section 1: Overview
	s1 := createSection(database, module.ID, "Overview of Load & Haul", 0)
	createSlide(database, s1.ID, "Load & Haul Operations", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Load & Haul Operations",
		"subtitle": "Understanding the engine room of the mining value chain",
		"bullets": []string{
			"How the site operates day to day",
			"Equipment roles and responsibilities",
			"Loading, hauling, and dumping fundamentals",
		},
	})
	createSlide(database, s1.ID, "Where Load & Haul Fits", "DIAGRAM", 1, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Load & Haul Cycle",
		"description": "Load & haul is the critical link between extracting material from the ground and delivering it to processing or permanent storage. The cycle — load, haul, dump, return — repeats continuously for every truck, every shift. It typically represents 40-60% of total mining operating costs, so small efficiency gains in any step of the cycle have outsized financial impact.",
		"diagramType": "image",
		"imageUrl":    "/assets/diagrams/load_haul_cycle.png",
	})
	createSlide(database, s1.ID, "Key Roles in Load & Haul", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Key Roles in Load & Haul",
		"introText": "Every person on the load & haul team has a critical role in keeping material moving safely and efficiently.",
		"items": []map[string]interface{}{
			{"title": "Excavator Operator", "description": "Loads trucks efficiently while maintaining dig face stability. Controls the pace of the entire operation."},
			{"title": "Truck Driver", "description": "Hauls material between loading and dumping points. Responsible for safe operation and maximising payload."},
			{"title": "Supervisor / Dispatcher", "description": "Allocates trucks to excavators, manages fleet priorities, and optimises cycle times."},
			{"title": "Spotter / Offsider", "description": "Guides trucks in confined areas, monitors dump edges, and assists with ground conditions."},
			{"title": "Road Maintenance Crew", "description": "Maintains haul road surface quality — directly impacts tyre life, fuel consumption, and cycle times."},
		},
	})
	createSlide(database, s1.ID, "Truck Spotting — Shovel Shadow Technique", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Truck Spotting — Shovel Shadow Technique",
		"introText": "Spotting is how the truck driver positions the tray under the excavator bucket. The shovel shadow method is the industry standard — consistent, safe, and maximises loading efficiency.",
		"items": []map[string]interface{}{
			{"title": "The Shadow Method", "description": "With the bucket held out at loading height, the excavator casts a shadow on the ground. The truck reverses into that shadow — when the shadow falls across the tray, the truck is in the bucket's swing arc. This gives the operator a clear, repeatable visual cue without radio guidance."},
			{"title": "Tyre Placement Target", "description": "Front tyre aligned to the bucket lip (visible through the rear-view mirror). This puts the tray centered under the bucket for a balanced load and minimum swing angle."},
			{"title": "Common Mistake — Under-Spotting", "description": "Truck stops short of the bucket. The excavator has to stretch and dump early, causing bucket stalling, partial loads, and uneven weight distribution."},
			{"title": "Common Mistake — Over-Spotting", "description": "Truck reverses too far — into the excavator's swing radius. This is a collision risk: the bucket or counterweight can strike the cab on the next swing. Always stop in the shadow, never past it."},
			{"title": "Single-Pass vs Two-Pass Loading", "description": "Single-pass = truck takes one load and leaves; fast cycle but requires perfect matching of bucket size to truck capacity. Two-pass (two trucks in sequence at the same face) = higher excavator utilisation on mismatched fleets, but doubles spotting operations and increases collision exposure."},
		},
	})
	createSlide(database, s1.ID, "Payload Management & On-Board Weighing", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Payload Management & On-Board Weighing",
		"introText": "Every tonne the truck doesn't carry — or carries over nominal — costs money. On-Board Weighing (OBW) is how the fleet keeps payload on target, load after load.",
		"items": []map[string]interface{}{
			{"title": "How OBW Works", "description": "Load cells or suspension pressure sensors measure the weight of material in the tray. Reading updates live as the bucket dumps. Operator sees current payload on in-cab display; supervisor sees fleet-wide payload distribution on FMS."},
			{"title": "Target Utilisation", "description": "Aim for >95% of nominal payload. A 290 t truck should consistently deliver 275–290 t loads. Loads below 275 t are 'light' — productivity lost. Loads above 290 t are overloads — damage and safety risk."},
			{"title": "Overloading Consequences", "description": "Tyre overheating and accelerated wear (tyres are ~$80k each on ultra-class), suspension strut damage, frame fatigue, haul road surface degradation from excess axle load. Most OEMs void warranty on chronic overloads."},
			{"title": "Under-Loading Consequences", "description": "Every light load = lost tonnes at the same fuel burn and cycle time. A 10% under-load across a fleet of 20 trucks = equivalent to parking 2 trucks. Cost per tonne goes up linearly."},
			{"title": "FMS Integration", "description": "OBW data streams into the FMS. Dashboards show payload distribution per truck, per operator, per shovel. Outlier operators get coaching; outlier shovels indicate bucket fill or face-presentation problems. OBW is the feedback loop that makes payload targets real."},
		},
	})
	createSlide(database, s1.ID, "FMS Dispatch — How Trucks Get Assigned", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "FMS Dispatch — How Trucks Get Assigned",
		"introText": "Modern open-cut mines don't direct trucks by radio. The Fleet Management System (FMS) assigns trucks to shovels and dumps in real time, based on live cycle data, queue lengths, and priority rules.",
		"items": []map[string]interface{}{
			{"title": "FMS Role", "description": "Real-time allocation of trucks to loading units (shovels, loaders) and dump destinations (crusher, waste dumps, stockpiles). The system continuously re-optimises assignments to minimise truck idle time and shovel idle time."},
			{"title": "In-Cab Display", "description": "Operators receive their next assignment on the in-cab display — truck ID, source, destination, route. Not radio-directed. This removes dispatcher-operator radio chatter and makes assignments auditable."},
			{"title": "Pre-Start Workflow", "description": "Operator logs in (card or PIN), enters shift and crew details, completes pre-start checks in the FMS. The system then assigns the first task — usually the nearest productive shovel with available queue capacity."},
			{"title": "Dynamic Reassignment", "description": "Mid-cycle reassignment happens constantly. If the assigned shovel's queue grows beyond threshold, the system may divert the truck to a closer shovel. If a crusher fills, ore trucks are re-routed to stockpile. Operators must follow the current assignment, not the previous one."},
			{"title": "When Dispatch Conflicts with Field Reality", "description": "If the assigned route is blocked, the dump is closed, or the ground is unsafe — do NOT force the task. Call the supervisor on radio, report the blocker, wait for reassignment. The FMS does not see everything on the ground; operators are the final safety gate."},
		},
	})
	createSlide(database, s1.ID, "Dozer Interaction at the Dump", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Dozer Interaction at the Dump",
		"introText": "The dozer is the partner of the haul truck at every dump. Its job is to maintain a safe, productive tipping face — and the truck-dozer interaction is one of the highest-risk interfaces on the mine.",
		"items": []map[string]interface{}{
			{"title": "Dozer Role at the Dump", "description": "Push tipped material over the dump edge (paddock dumps) or down the highwall slope. Rebuild the windrow after each tip. Push fines back from the travel surface. Maintain a flat, well-drained tipping pad for trucks."},
			{"title": "Windrow — The Safety Berm", "description": "Continuous earth berm along the dump edge. Minimum height = 0.5 × largest truck tyre height (typically 1.5–2 m for ultra-class trucks). The windrow is the last physical barrier preventing a truck from going over the edge. Never dump on a section with a broken or missing windrow."},
			{"title": "Reverse Distance Rule", "description": "Truck rear tyres must not come closer than 1 m from the windrow. The truck reverses until rear tyres contact the windrow, then stops. Reversing into the windrow at speed or pushing through it is a fatal-class risk event."},
			{"title": "No Simultaneous Operation", "description": "Dozer and truck must not operate at the same dump face at the same time. Clear radio communication: dozer announces 'clear' before the next truck can approach; truck announces 'tipping' before the dozer re-enters. One machine at a time on the tipping pad."},
			{"title": "Dozer Cleanup", "description": "Between truck cycles, the dozer pushes fines and spill material back from the haul road onto the tipping face, and re-levels the pad. This maintains the surface for the next truck and prevents tyre damage from loose rock on the travel surface."},
		},
	})
	createSlide(database, s1.ID, "Mining Value Chain Overview", "PROCESS", 7, map[string]interface{}{
		"type":        "process",
		"heading":     "Mining Value Chain Overview",
		"description": "The full lifecycle of material from in-situ rock to final product.",
		"steps": []map[string]interface{}{
			{"title": "Exploration & Planning", "description": "Geological surveys, resource modelling, mine planning and scheduling."},
			{"title": "Drill & Blast", "description": "Drill patterns designed, holes charged with explosives, rock fragmented to target size."},
			{"title": "Load", "description": "Excavators or loaders fill haul trucks. Material classified as ore or waste based on grade control."},
			{"title": "Haul", "description": "Trucks transport material via haul roads to designated dump points, crusher, or ROM pad."},
			{"title": "Dump & Process", "description": "Material tipped at destination. Ore goes to crusher/processing; waste to designated dumps."},
			{"title": "Rehabilitation", "description": "Progressive rehab of dump sites and mined-out areas to meet environmental commitments."},
		},
	})

	// Section 2: Loading Operations
	s2 := createSection(database, module.ID, "Loading Operations", 1)
	createSlide(database, s2.ID, "Loading Area Fundamentals", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Loading Area Fundamentals",
		"introText": "The loading area (dig face) is where value creation begins. A well-managed loading area maximises productivity and minimises risk.",
		"items": []map[string]interface{}{
			{"title": "Bench Height", "description": "Should be 1.0 to 1.5 times the machine height for optimal digging efficiency and safety. Too high risks face collapse; too low reduces bucket fill.", "imageUrl": "https://assets.miningacademy.com/images/bench-height-excavator.jpg"},
			{"title": "Face Angle", "description": "Maintained at design angle (typically 60-75°) to prevent slips and provide stable working conditions."},
			{"title": "Floor Conditions", "description": "Must be level, well-drained, and compacted. Poor floors cause instability, tyre damage, and reduced productivity."},
			{"title": "Truck Presentation Area", "description": "Sufficient space for trucks to enter, position, and exit safely. Minimum 2x truck width clearance."},
			{"title": "Escape Routes", "description": "Always maintain at least two exit paths from the loading area for emergency evacuation."},
		},
	})
	createSlide(database, s2.ID, "Excavator Types", "COMPARISON", 1, map[string]interface{}{
		"type":        "comparison",
		"heading":     "Excavator Types in Mining",
		"description": "Different excavator configurations suit different mining conditions.",
		"items": []map[string]interface{}{
			{
				"name":        "Hydraulic Excavator (Backhoe)",
				"description": "Most common in mining. Digs below track level with a pulling motion.",
				"pros":        []string{"Excellent digging force", "Good reach and depth", "Versatile — can work in confined spaces", "Lower floor preparation needed"},
				"cons":        []string{"Slower cycle time than rope shovels", "Higher maintenance cost per tonne", "Limited bucket size compared to rope shovels"},
				"specs":       []map[string]string{{"label": "Bucket Size", "value": "15-45 m³"}, {"label": "Max Dig Depth", "value": "8-12m"}, {"label": "Cycle Time", "value": "25-35 sec"}},
				"imageUrl":    "https://assets.miningacademy.com/images/excavator-hydraulic-backhoe.jpg",
			},
			{
				"name":        "Rope Shovel (Electric)",
				"description": "Purpose-built for high-volume production loading in large open pits.",
				"pros":        []string{"Highest productivity per hour", "Lower cost per tonne at scale", "Extremely long service life (30+ years)", "Electric drive — no diesel costs"},
				"cons":        []string{"Very high capital cost", "Requires flat, prepared floor", "Limited mobility — needs ramp access", "Requires trailing cable management"},
				"specs":       []map[string]string{{"label": "Bucket Size", "value": "40-80 m³"}, {"label": "Max Dig Height", "value": "15-18m"}, {"label": "Cycle Time", "value": "28-35 sec"}},
				"imageUrl":    "https://assets.miningacademy.com/images/excavator-rope-shovel.jpg",
			},
			{
				"name":        "Front-End Loader",
				"description": "Rubber-tyred loader used for supplementary loading, clean-up, and stockpile work.",
				"pros":        []string{"Highly mobile — can self-relocate", "Good for multiple short-term faces", "Lower capital cost", "Can load from stockpiles"},
				"cons":        []string{"Lower digging force than excavators", "Higher cost per tonne for primary loading", "Tyre wear is significant cost", "Requires more floor preparation"},
				"specs":       []map[string]string{{"label": "Bucket Size", "value": "10-25 m³"}, {"label": "Payload", "value": "20-45 tonnes"}, {"label": "Cycle Time", "value": "30-45 sec"}},
				"imageUrl":    "https://assets.miningacademy.com/images/excavator-front-end-loader.jpg",
			},
		},
	})
	createSlide(database, s2.ID, "Loading Techniques", "PROCESS", 2, map[string]interface{}{
		"type":        "process",
		"heading":     "Optimal Loading Sequence",
		"description": "The standard loading sequence for maximising productivity and safety.",
		"steps": []map[string]interface{}{
			{"title": "Truck Approach", "description": "Truck enters the loading area following the designated access path. Driver signals arrival to the excavator operator.", "warnings": []string{"Maintain safe following distance from the truck ahead"}},
			{"title": "Truck Positioning", "description": "Truck reverses to the excavator's loading side, positioning the tray directly under the bucket path. Aim for minimum swing angle."},
			{"title": "Loading Passes", "description": "Excavator fills truck in 3-5 passes (optimal is 4). Each pass targets even distribution across the tray."},
			{"title": "Load Check", "description": "Final pass levels the load. Operator visually confirms payload is within target range and evenly distributed."},
			{"title": "Truck Departure", "description": "Driver receives clear signal, exits the loading area via the designated route. Next truck moves into position immediately."},
		},
	})

	// Section 3: Hauling
	s3 := createSection(database, module.ID, "Hauling Operations", 2)
	createSlide(database, s3.ID, "Scale of a Mining Haul Truck", "IMAGE", 0, map[string]interface{}{
		"type":        "image",
		"heading":     "Scale of a Mining Haul Truck",
		"description": "The ultra-class haul trucks used in large open-cut operations are among the largest land vehicles on Earth. Standing next to one for the first time is a disorienting experience — the rear tyre alone stands taller than a person.",
		"imageUrl":    "https://assets.miningacademy.com/images/haul-truck-scale-person.jpg",
	})
	createSlide(database, s3.ID, "Haul Truck Types", "EQUIPMENT", 1, map[string]interface{}{
		"type":        "equipment",
		"heading":     "Ultra-Class Haul Truck",
		"description": "The CAT 797F is one of the largest mechanical-drive haul trucks in the world, designed for high-production mining.",
		"specifications": []map[string]string{
			{"label": "Payload Capacity", "value": "400 short tons (363 tonnes)"},
			{"label": "Gross Machine Weight", "value": "687.5 tonnes"},
			{"label": "Engine Power", "value": "4,000 HP (2,983 kW)"},
			{"label": "Top Speed (loaded)", "value": "42 km/h"},
			{"label": "Body Volume", "value": "220-310 m³"},
			{"label": "Fuel Tank", "value": "6,800 litres"},
			{"label": "Tyre Size", "value": "59/80R63"},
			{"label": "Tyre Cost", "value": "~$80,000 each"},
		},
		"advantages": []string{
			"Lowest cost per tonne for long-haul distances",
			"Matched to large excavators (4-pass loading)",
			"Advanced payload monitoring systems",
			"Autonomous haulage system (AHS) capable",
		},
		"limitations": []string{
			"Extremely high capital cost ($8-12M per unit)",
			"Requires wide, well-maintained haul roads",
			"Long lead time for tyre replacement",
			"Significant infrastructure requirements",
		},
	})
	createSlide(database, s3.ID, "Haul Road Design", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Haul Road Design Principles",
		"introText": "Haul roads are the arteries of the mine. Road quality directly impacts safety, tyre life, fuel consumption, and productivity.",
		"items": []map[string]interface{}{
			{"title": "Width", "description": "Minimum 3.5x the widest truck for two-way traffic. Single lane minimum 2x truck width plus safety berm."},
			{"title": "Gradient", "description": "Maximum sustained grade of 10% (1 in 10). Ramps typically designed at 8%. Steeper grades dramatically increase fuel consumption and brake wear."},
			{"title": "Super-elevation", "description": "Curves banked at 2-4% to counteract centrifugal force. Reduces tyre wear and improves stability at speed."},
			{"title": "Surface Material", "description": "Well-graded crushed rock (wearing course) over compacted base. Surface maintained with graders and water carts."},
			{"title": "Drainage", "description": "Cross-fall of 2-3% to shed water. Drains and culverts at low points. Standing water destroys road surface and is a safety hazard."},
			{"title": "Safety Berms", "description": "Minimum berm height = axle height of largest truck. Continuous berms on all elevated sections and dump edges.", "imageUrl": "https://assets.miningacademy.com/images/haul-road-safety-berm.jpg"},
		},
	})
	createSlide(database, s3.ID, "Cycle Time Components", "INTERACTIVE", 3, map[string]interface{}{
		"type":            "interactive",
		"heading":         "Haul Truck Cycle Time Calculator",
		"description":     "Adjust parameters to see how they affect total cycle time and productivity.",
		"interactionType": "calculator",
		"controls": []map[string]interface{}{
			{"id": "loadTime", "label": "Loading Time (min)", "type": "range", "min": 2, "max": 8, "step": 0.5, "default": 4},
			{"id": "haulDist", "label": "Haul Distance (km)", "type": "range", "min": 1, "max": 10, "step": 0.5, "default": 3},
			{"id": "haulSpeed", "label": "Avg Haul Speed (km/h)", "type": "range", "min": 15, "max": 45, "step": 5, "default": 30},
			{"id": "dumpTime", "label": "Dump Time (min)", "type": "range", "min": 1, "max": 5, "step": 0.5, "default": 2},
			{"id": "returnSpeed", "label": "Avg Return Speed (km/h)", "type": "range", "min": 20, "max": 50, "step": 5, "default": 35},
		},
		"outputs": []map[string]interface{}{
			{"id": "haulTime", "label": "Haul Time", "formula": "(haulDist / haulSpeed) * 60", "unit": "min"},
			{"id": "returnTime", "label": "Return Time", "formula": "(haulDist / returnSpeed) * 60", "unit": "min"},
			{"id": "totalCycle", "label": "Total Cycle Time", "formula": "loadTime + (haulDist / haulSpeed) * 60 + dumpTime + (haulDist / returnSpeed) * 60", "unit": "min"},
			{"id": "cyclesPerHour", "label": "Cycles Per Hour", "formula": "60 / (loadTime + (haulDist / haulSpeed) * 60 + dumpTime + (haulDist / returnSpeed) * 60)"},
		},
		"feedback": []map[string]interface{}{
			{"condition": "loadTime + (haulDist / haulSpeed) * 60 + dumpTime + (haulDist / returnSpeed) * 60 > 30", "message": "Cycle time exceeds 30 minutes — consider adding trucks or reducing haul distance.", "type": "warning"},
			{"condition": "loadTime + (haulDist / haulSpeed) * 60 + dumpTime + (haulDist / returnSpeed) * 60 < 15", "message": "Excellent cycle time — this is a highly productive configuration.", "type": "success"},
		},
	})

	// Section 4: Dumping Operations
	s4 := createSection(database, module.ID, "Dumping Operations", 3)
	createSlide(database, s4.ID, "Types of Dumps", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Types of Dump Points",
		"introText": "Material is directed to different dump types based on classification (ore vs waste) and operational requirements.",
		"items": []map[string]interface{}{
			{"title": "Waste Dump (Highwall)", "description": "Permanent storage for waste rock. Built in lifts with engineered design angles. Progressive rehabilitation required."},
			{"title": "Paddock Dump", "description": "Flat dump area where trucks tip in rows. Used when vertical dump heights are restricted or for temporary storage."},
			{"title": "Crusher Dump", "description": "Trucks tip directly into the primary crusher. Highest value destination — direct feed to processing plant."},
			{"title": "ROM Pad", "description": "Run-of-mine stockpile. Temporary storage between mining and processing. Managed by grade for blending."},
			{"title": "Road Maintenance Dump", "description": "Selected waste material used for road construction and repair. Specific material specs required."},
			{"title": "Backfill", "description": "Waste material returned to mined-out areas of the pit. Reduces long-term rehabilitation liability."},
		},
	})
	createSlide(database, s4.ID, "Safe Dumping Procedures", "PROCESS", 1, map[string]interface{}{
		"type":        "process",
		"heading":     "Safe Dumping Procedure",
		"description": "Follow this procedure every time. No exceptions.",
		"steps": []map[string]interface{}{
			{"title": "Approach", "description": "Approach the dump edge slowly in first gear. Observe all signage and spotter instructions.", "warnings": []string{"Never approach a dump edge at speed"}},
			{"title": "Position", "description": "Align truck perpendicular to the dump crest. Stop well back from the edge — minimum 5m from crest.", "warnings": []string{"Always check ground conditions before reversing to edge"}},
			{"title": "Reverse", "description": "Reverse slowly toward the berm/windrow. Spotter guides from a safe position, visible to the driver.", "warnings": []string{"Stop immediately if you lose sight of the spotter"}},
			{"title": "Contact Berm", "description": "Reverse until rear tyres contact the safety berm. Come to a complete stop. Apply park brake."},
			{"title": "Tip", "description": "Raise body slowly. If load sticks, do NOT drive forward with body up. Lower body and seek assistance."},
			{"title": "Lower & Depart", "description": "Lower body completely before moving. Check mirrors. Depart via designated exit route."},
		},
	})

	// Section 5: Safety
	s5 := createSection(database, module.ID, "Safety in Load & Haul", 4)
	createSlide(database, s5.ID, "Critical Safety Rules", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Critical Safety Rules",
		"introText": "These rules are non-negotiable. Violations may result in immediate removal from site.",
		"items": []map[string]interface{}{
			{"title": "Seatbelts", "description": "Worn at ALL times when the vehicle is in motion. No exceptions."},
			{"title": "Right of Way", "description": "Loaded trucks ALWAYS have right of way over empty trucks and light vehicles. Yield at all intersections."},
			{"title": "Speed Limits", "description": "Strictly enforced. Haul roads typically 40-60 km/h. Ramps 15-25 km/h. Near people 15 km/h max."},
			{"title": "Exclusion Zones", "description": "Never enter an excavator's swing radius without direct communication and confirmation from the operator."},
			{"title": "Fatigue Management", "description": "Report fatigue immediately. Use fatigue detection systems. Take breaks as scheduled."},
			{"title": "Pre-Start Checks", "description": "Complete a thorough pre-start inspection every shift. Report all defects. Do not operate defective equipment."},
		},
	})
	createSlide(database, s5.ID, "Hazard Awareness", "DIAGRAM", 1, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Common Load & Haul Hazards",
		"description": "Understanding hazards is the first step to controlling risk.",
		"diagramType": "process",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"label": "Collision", "description": "Truck-to-truck, truck-to-equipment", "highlight": true},
				{"label": "Rollover", "description": "Berm failure, soft ground, speed", "highlight": true},
				{"label": "Falls", "description": "Dump edge failure, unstable ground", "highlight": true},
				{"label": "Tyre Failure", "description": "Blowout, fire, projectile risk", "highlight": false},
				{"label": "Fatigue", "description": "Reduced alertness, micro-sleeps", "highlight": false},
			},
		},
	})

	createSlide(database, s5.ID, "What a Berm Failure Looks Like", "IMAGE", 2, map[string]interface{}{
		"type":        "image",
		"heading":     "What a Berm Failure Looks Like",
		"description": "An ultra-class haul truck weighs up to 560 tonnes when loaded. When a berm fails or is absent, there is no effective recovery — the truck goes over the edge. This is not an abstract risk: berm failures at dump edges are among the leading causes of serious injuries and fatalities in open-cut mining.",
		"imageUrl":    "https://assets.miningacademy.com/images/haul-truck-rollover-berm-failure.jpg",
	})

	// Section 6: Efficiency & Optimization
	s6 := createSection(database, module.ID, "Efficiency & Optimization", 5)
	createSlide(database, s6.ID, "Truck-Excavator Matching", "INTERACTIVE", 0, map[string]interface{}{
		"type":            "interactive",
		"heading":         "Truck-Excavator Match Factor",
		"description":     "The match factor determines how well the truck fleet is matched to the excavator. Target: 0.85-1.15",
		"interactionType": "calculator",
		"controls": []map[string]interface{}{
			{"id": "numTrucks", "label": "Number of Trucks", "type": "range", "min": 2, "max": 10, "step": 1, "default": 5},
			{"id": "truckCycle", "label": "Truck Cycle Time (min)", "type": "range", "min": 10, "max": 40, "step": 1, "default": 22},
			{"id": "loadTime2", "label": "Loading Time (min)", "type": "range", "min": 2, "max": 8, "step": 0.5, "default": 4},
		},
		"outputs": []map[string]interface{}{
			{"id": "matchFactor", "label": "Match Factor", "formula": "(numTrucks * loadTime2) / truckCycle"},
			{"id": "trucksNeeded", "label": "Ideal Truck Count", "formula": "Math.round(truckCycle / loadTime2)"},
		},
		"feedback": []map[string]interface{}{
			{"condition": "(numTrucks * loadTime2) / truckCycle < 0.85", "message": "Under-trucked — the excavator will have idle time waiting for trucks. Consider adding trucks.", "type": "warning"},
			{"condition": "(numTrucks * loadTime2) / truckCycle > 1.15", "message": "Over-trucked — trucks are queuing at the excavator. Consider removing trucks or opening another face.", "type": "warning"},
			{"condition": "(numTrucks * loadTime2) / truckCycle >= 0.85 && (numTrucks * loadTime2) / truckCycle <= 1.15", "message": "Good match — the fleet is well balanced. Minimal truck queuing and excavator idle time.", "type": "success"},
		},
	})
	createSlide(database, s6.ID, "Key Performance Indicators", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Load & Haul KPIs",
		"introText": "These are the metrics that drive load & haul performance.",
		"items": []map[string]interface{}{
			{"title": "Tonnes per Operating Hour", "description": "Primary productivity measure. Total tonnes moved divided by operating hours. Target varies by equipment class."},
			{"title": "Cost per Tonne", "description": "Total load & haul cost divided by tonnes moved. Includes fuel, tyres, maintenance, labour, and depreciation."},
			{"title": "Truck Utilisation", "description": "Percentage of available time the truck is actually hauling. Target: >85%. Losses from queuing, delays, maintenance."},
			{"title": "Average Payload", "description": "Average tonnes per truck load. Monitored against target to prevent under-loading (lost productivity) and overloading (damage/safety)."},
			{"title": "Fuel Burn Rate", "description": "Litres per hour or litres per tonne-kilometre. Affected by road grade, speed, payload, and operator technique."},
		},
	})

	// Section 7: Knowledge Check
	s7 := createSection(database, module.ID, "Knowledge Check", 6)
	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Load & Haul Operations Assessment",
		PassMark:         80,
		MaxAttempts:      3,
		ShuffleQuestions: true,
		ShuffleOptions:   true,
		ShowResults:      true,
		SectionID:        &s7.ID,
		ModuleID:         &module.ID,
	}
	database.Create(&quiz)

	questions := []struct {
		question    string
		explanation string
		options     []struct {
			text      string
			isCorrect bool
		}
	}{
		{
			question:    "What percentage of total mining operating costs does load & haul typically represent?",
			explanation: "Load & haul operations typically account for 40-60% of total mining operating costs, making efficiency improvements in this area extremely impactful.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"10-20%", false},
				{"20-30%", false},
				{"40-60%", true},
				{"70-80%", false},
			},
		},
		{
			question:    "What is the ideal bench height relative to the excavator?",
			explanation: "The bench height should be 1.0 to 1.5 times the machine height for optimal digging efficiency and safety.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Half the stick length", false},
				{"1.0 to 1.5 times the machine height", true},
				{"Twice the bucket width", false},
				{"Equal to the truck tray height", false},
			},
		},
		{
			question:    "What is the optimal number of loading passes for a haul truck?",
			explanation: "3-5 passes is the target, with 4 being optimal. This balances excavator efficiency (not too many small passes) with truck utilisation (not waiting too long).",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"1-2 passes", false},
				{"3-5 passes (4 optimal)", true},
				{"6-8 passes", false},
				{"As many as needed", false},
			},
		},
		{
			question:    "What is the maximum sustained grade for a haul road?",
			explanation: "Maximum sustained grade is typically 10% (1 in 10), with ramps designed at 8%. Steeper grades dramatically increase fuel consumption and brake wear.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"5%", false},
				{"10%", true},
				{"15%", false},
				{"20%", false},
			},
		},
		{
			question:    "Which vehicle always has right of way on a mine site?",
			explanation: "Loaded trucks always have right of way because they are heavier, harder to stop, and represent the highest-value activity on the haul road.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Light vehicles", false},
				{"Empty trucks", false},
				{"Loaded trucks", true},
				{"Graders", false},
			},
		},
		{
			question:    "What is the minimum haul road width for two-way traffic?",
			explanation: "Minimum 3.5 times the width of the widest truck for two-way traffic, providing adequate passing clearance.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"2x truck width", false},
				{"3x truck width", false},
				{"3.5x truck width", true},
				{"4x truck width", false},
			},
		},
		{
			question:    "What is the target match factor range for a truck-excavator fleet?",
			explanation: "A match factor of 0.85-1.15 indicates the fleet is well balanced with minimal truck queuing and excavator idle time.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"0.50-0.75", false},
				{"0.85-1.15", true},
				{"1.20-1.50", false},
				{"1.50-2.00", false},
			},
		},
		{
			question:    "What is the minimum safe distance from the dump crest before reversing?",
			explanation: "Stop a minimum of 5 metres from the crest before beginning the reversing sequence. This provides a safety buffer in case of ground failure.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"1 metre", false},
				{"3 metres", false},
				{"5 metres", true},
				{"10 metres", false},
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

	// Knowledge check slide linking to quiz
	createSlide(database, s7.ID, "Load & Haul Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Load & Haul Operations Assessment",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{"question": "What percentage of mining costs does load & haul represent?", "options": []string{"10-20%", "20-30%", "40-60%", "70-80%"}, "correctIndex": 2, "explanation": "Load & haul = 40-60% of total mining operating costs."},
			{"question": "What is the ideal bench height?", "options": []string{"Half stick length", "1.0-1.5x machine height", "Twice bucket width", "Equal to tray height"}, "correctIndex": 1, "explanation": "1.0 to 1.5 times machine height for optimal efficiency and safety."},
			{"question": "Optimal loading passes for a haul truck?", "options": []string{"1-2", "3-5 (4 optimal)", "6-8", "As many as needed"}, "correctIndex": 1, "explanation": "4 passes balances excavator efficiency with truck utilisation."},
		},
	})

	// Completion slide
	createSlide(database, s7.ID, "Module Complete", "COMPLETION", 1, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "Congratulations! You have completed the Load & Haul Operations module.",
		"keyTakeaways": []string{
			"Load & haul represents 40-60% of mining costs — efficiency here has the biggest impact",
			"Bench height, floor conditions, and truck presentation are critical to loading productivity",
			"Haul road design directly affects tyre life, fuel consumption, and safety",
			"Match factor (0.85-1.15) balances truck and excavator utilisation",
			"Safety rules around dump edges, exclusion zones, and right of way are non-negotiable",
		},
		"showCertificate": true,
	})

	log.Println("Load & Haul module seeded with 7 sections and quiz")
}

func createSection(db *gorm.DB, moduleID, title string, sortOrder int) models.Section {
	s := models.Section{
		ID:        uuid.New().String(),
		Title:     title,
		SortOrder: sortOrder,
		ModuleID:  moduleID,
	}
	db.Create(&s)
	return s
}

func createSlide(db *gorm.DB, sectionID, title string, slideType string, sortOrder int, content map[string]interface{}) {
	contentJSON, _ := json.Marshal(content)
	db.Create(&models.Slide{
		ID:        uuid.New().String(),
		Title:     title,
		Type:      models.SlideType(slideType),
		Content:   contentJSON,
		SortOrder: sortOrder,
		SectionID: sectionID,
	})
}
