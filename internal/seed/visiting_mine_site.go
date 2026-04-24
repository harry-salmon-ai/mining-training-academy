package seed

import (
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedVisitingMineSiteModule(database *gorm.DB, force bool) {
	var existing models.Module
	database.Where("slug = ?", "visiting-a-mine-site").First(&existing)
	if existing.ID != "" {
		if force {
			log.Println("SEED=true: deleting existing Visiting a Mine Site module for re-seed...")
			forceDeleteModule(database, "visiting-a-mine-site")
			log.Println("Existing Visiting a Mine Site module deleted")
		} else {
			log.Println("Visiting a Mine Site module already exists, skipping")
			return
		}
	}

	log.Println("Seeding Visiting a Mine Site module...")

	// Ensure "Site Orientation" category exists
	var category models.Category
	database.Where("slug = ?", "site-orientation").First(&category)
	if category.ID == "" {
		category = models.Category{
			ID:          uuid.New().String(),
			Name:        "Site Orientation",
			Slug:        "site-orientation",
			Color:       "#0D9488",
			Icon:        "O",
			SortOrder:   7,
			IsActive:    true,
			Description: ptrStr("Preparing engineers for their first mine site experience"),
		}
		database.Create(&category)
		log.Println("Created Site Orientation category")
	}

	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed Visiting a Mine Site: category or admin not found")
		return
	}

	desc := "Everything a software engineer needs to know before stepping onto an open pit mine site for the first time — etiquette, departments, who to talk to about technology, and the lingo that separates a miner from an outsider."
	dur := 60
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Visiting a Mine Site",
		Slug:        "visiting-a-mine-site",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelFoundation,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"site visit", "etiquette", "mine lingo", "departments", "value chain", "PPE", "safety", "AHS", "FMS", "open pit", "orientation"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 1: Mine Etiquette ───────────────────────────────────────────
	s1 := createSection(database, module.ID, "Mine Etiquette", 0)

	createSlide(database, s1.ID, "Mine Etiquette", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Mine Etiquette",
		"subtitle": "The rules that keep you safe and welcomed on site",
		"bullets": []string{
			"PPE requirements — what's mandatory, what's site-specific",
			"How to behave on haul roads, around equipment, and in blast zones",
			"Radio protocol and communication on site",
		},
	})

	createSlide(database, s1.ID, "Personal Protective Equipment", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Personal Protective Equipment (PPE)",
		"introText": "PPE is non-negotiable. You will not be allowed past the gate without it. Every site has a minimum standard — some add extras depending on the operation.",
		"items": []map[string]interface{}{
			{"title": "Hard Hat", "description": "Worn at all times outside a vehicle or building. Must meet AS/NZS 1801 or equivalent. Colour often indicates role — white for visitors, blue for operators, red for supervisors. Check with your host."},
			{"title": "High-Visibility Clothing", "description": "Long-sleeve hi-vis shirt and pants, or hi-vis vest over long sleeves. Must have reflective striping for night visibility. No exceptions, even on a quick walk."},
			{"title": "Steel-Cap Boots", "description": "Steel toe and metatarsal protection. Must be lace-up, above the ankle. No slip-ons, no sneakers, no exceptions."},
			{"title": "Safety Glasses", "description": "Clear or tinted depending on conditions. Worn whenever outside. Prescription safety glasses are available — organise in advance."},
			{"title": "Hearing Protection", "description": "Earplugs or earmuffs required near operating equipment, processing plants, and blast zones. Carry them at all times — you won't always get warning."},
			{"title": "Site-Specific Extras", "description": "Some sites require dust masks, gloves, or fall protection in specific areas. Your site induction will cover these. If in doubt, ask your host before walking in."},
		},
	})

	createSlide(database, s1.ID, "Site Access and Sign-In", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Site Access and Sign-In",
		"introText": "Every mine site tracks who is on site at all times. This is a legal requirement and a safety one — if something goes wrong, they need to know who to account for.",
		"items": []map[string]interface{}{
			{"title": "Visitor Induction", "description": "Before your first visit, you'll complete a site-specific induction. This may be online (before arrival) or in person (at the gate). It covers emergency procedures, site rules, and hazard awareness. Take it seriously — people have been turned away for not completing it."},
			{"title": "Sign-In / Sign-Out", "description": "You must sign in on arrival and sign out on departure. This is usually at the gatehouse or security office. If you don't sign out, someone will come looking for you."},
			{"title": "Escort Requirements", "description": "Most sites require visitors to be escorted at all times by an inducted employee. You are their responsibility. Stay with your escort."},
			{"title": "Vehicle Registration", "description": "If you're driving a vehicle onto site, it must be registered, inspected (pre-start check), and fitted with a flag, flashing light, and UHF radio."},
		},
	})

	createSlide(database, s1.ID, "Light Vehicle and Haul Road Rules", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Light Vehicle and Haul Road Rules",
		"introText": "A loaded haul truck weighs up to 600 tonnes and cannot stop quickly. You are invisible to it. These rules exist because light vehicle interactions with haul trucks are one of the leading causes of fatalities in mining.",
		"items": []map[string]interface{}{
			{"title": "Give Way — Always", "description": "Haul trucks have right of way. Period. If you see one, pull over well clear of the road and wait. Do not try to overtake, squeeze past, or assume the driver has seen you."},
			{"title": "Speed Limits", "description": "Strictly enforced. Typically 60 km/h on main haul roads, 40 km/h near operations, 15 km/h near people. Your vehicle will likely have GPS speed monitoring."},
			{"title": "Headlights On", "description": "Headlights and rotating beacon on at all times, day and night. This is how truck operators spot you."},
			{"title": "No Phones While Driving", "description": "Zero tolerance. Not hands-free, not for maps, not for a quick check. Pull over and stop if you need your phone."},
			{"title": "Seatbelts", "description": "Worn at all times when the vehicle is in motion. Every occupant. Front and back."},
			{"title": "Parking", "description": "Park in designated areas only. Never park on a haul road, near active equipment, or past a berm. Wheels turned into the berm where possible."},
		},
	})

	createSlide(database, s1.ID, "Radio Protocol", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Radio Protocol",
		"introText": "UHF radio is the primary communication tool on a mine site. It's not like a phone call — everyone on the channel hears everything. Keep it short, clear, and professional.",
		"items": []map[string]interface{}{
			{"title": "Channel Discipline", "description": "Each area of the mine has a designated channel. The control room monitors the main channel. Switch to the correct channel before transmitting. Listen before you speak — don't talk over someone."},
			{"title": "Calling Your Position", "description": "When entering a new area or haul road, call your vehicle ID, location, and direction. Example: 'Light vehicle 42, entering Ramp 3, heading south.' This tells truck operators you're there."},
			{"title": "Keep It Brief", "description": "Radio is shared by dozens of operators. Say what you need to say and get off. No chatter, no jokes, no long explanations."},
			{"title": "Emergency Calls", "description": "If you hear 'Mayday' or 'Emergency' on the radio, stop transmitting immediately and listen. Do not use the channel until cleared."},
			{"title": "If You Don't Have a Radio", "description": "You should not be driving on a mine site without one. If your vehicle doesn't have a working UHF, do not proceed. Tell your escort."},
		},
	})

	createSlide(database, s1.ID, "Blast Protocols", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "Blast Protocols",
		"introText": "Open pit mines blast regularly — often daily. Blasting involves explosives, exclusion zones, and strict timing. If you're on site during a blast, you need to know what to do.",
		"items": []map[string]interface{}{
			{"title": "The Blast Schedule", "description": "Most sites blast at a set time each day (commonly around midday or shift change). Your host should tell you the schedule. If not, ask."},
			{"title": "Exclusion Zones", "description": "Before a blast, the area is cleared and exclusion zones are established. These are marked with signs, barriers, and sentries. Do not enter an exclusion zone under any circumstances."},
			{"title": "The Siren Sequence", "description": "A series of sirens signals the blast: long siren = area is being cleared; short repeated sirens = blast is imminent; all clear siren = safe to return. The exact pattern varies by site — learn it during induction."},
			{"title": "What to Do", "description": "When you hear the blast siren, go to the designated muster point. Stay there until the all clear. Don't try to watch the blast from a 'good spot' — flyrock is unpredictable and lethal."},
		},
	})

	createSlide(database, s1.ID, "No-Go Zones and General Behaviour", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "No-Go Zones and General Behaviour",
		"introText": "Mine sites are full of areas that look accessible but aren't. Respect the boundaries — they exist because someone was hurt or killed when they didn't.",
		"items": []map[string]interface{}{
			{"title": "Berms", "description": "A berm is a mound of material along the edge of a road, bench, or dump. It's a physical barrier to stop vehicles going over the edge. Never drive past a berm. Never park with your front wheels past a berm. If a berm is missing or damaged, report it immediately."},
			{"title": "Bench Edges", "description": "The edge of a mining bench can be unstable. Stay well back. If there's no berm, don't approach."},
			{"title": "Active Equipment Zones", "description": "Do not approach operating excavators, drills, or dozers without radio contact and confirmation from the operator. Their blind spots are enormous."},
			{"title": "The Buddy System", "description": "In remote areas of the pit, don't go alone. Always travel with someone who can raise the alarm if something goes wrong."},
			{"title": "Muster Points", "description": "Know where your nearest muster point is at all times. In an emergency, go there and wait to be accounted for."},
			{"title": "Respect the Site", "description": "No alcohol, no drugs, no smoking outside designated areas. No photography without permission. Follow instructions from any site personnel immediately, even if you outrank them in the office."},
		},
	})

	// ── Section 2: The Mining Value Chain & Departments ─────────────────────
	s2 := createSection(database, module.ID, "The Mining Value Chain & Departments", 1)

	createSlide(database, s2.ID, "The Mining Value Chain", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "The Mining Value Chain & Departments",
		"subtitle": "Who does what, and where they sit in the flow from rock to product",
		"bullets": []string{
			"The end-to-end open pit mining value chain",
			"Each department's role and what they care about",
			"How departments relate to the technology stack",
		},
	})

	createSlide(database, s2.ID, "The Open Pit Value Chain", "PROCESS", 1, map[string]interface{}{
		"type":        "process",
		"heading":     "The Open Pit Mining Value Chain",
		"description": "Every open pit mine follows this sequence. Each step has a dedicated department, and each department has its own technology needs. Understanding this chain is the foundation for knowing who to talk to.",
		"steps": []map[string]interface{}{
			{"title": "Exploration & Geology", "description": "Find the ore body, characterise it, estimate the resource."},
			{"title": "Mine Planning", "description": "Design the pit, schedule extraction, balance ore and waste."},
			{"title": "Drill & Blast", "description": "Drill holes, load explosives, fragment rock to target size."},
			{"title": "Load & Haul", "description": "Excavators load trucks, trucks haul to crusher or waste dump."},
			{"title": "Processing", "description": "Crush, grind, extract. Turn ore into saleable product."},
			{"title": "Rehabilitation", "description": "Restore mined land. Progressive and ongoing."},
		},
	})

	createSlide(database, s2.ID, "Geology & Resource", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Geology & Resource",
		"introText": "The origin of all mine data. Geologists define where the ore is, what grade it is, and how much exists. They build the block model that every other department depends on.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Resource estimation, geological mapping, core logging, grade control sampling. They decide what's ore and what's waste."},
			{"title": "Technologies They Use", "description": "Geological modelling software (Leapfrog, Vulcan, Surpac), drill hole databases, assay management systems."},
			{"title": "When You'd Talk to Them", "description": "Rarely, unless your software touches geological data, grade control, or ore-body modelling. They're typically in an office on site, not in the pit."},
		},
	})

	createSlide(database, s2.ID, "Mine Planning & Technical Services", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Mine Planning & Technical Services",
		"introText": "The architects of the operation. They design the pit, schedule what gets mined and when, and reconcile what was planned against what actually happened.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Long-term mine plans (life of mine), short-term plans (weekly/daily dig plans), survey, grade control, and reconciliation. They produce the dig plan that drives everything downstream."},
			{"title": "Technologies They Use", "description": "Mine planning software (Deswik, MineSched, Whittle), survey equipment (GPS, drones, total stations), reconciliation tools."},
			{"title": "When You'd Talk to Them", "description": "If your software uses planned vs actual data, scheduling data, spatial/survey data, or anything related to where and when material is mined."},
		},
	})

	createSlide(database, s2.ID, "Drill & Blast", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Drill & Blast",
		"introText": "They break the rock. Everything upstream is planning; everything downstream depends on them getting the fragmentation right.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Design blast patterns, drill holes to spec, load explosives, fire blasts, measure fragmentation outcomes. They control when and where material gets broken."},
			{"title": "Technologies They Use", "description": "Drill navigation/guidance systems, blast design software (SHOTPlus, BlastLogic), MWD (Measure While Drilling) systems, fragmentation analysis tools."},
			{"title": "When You'd Talk to Them", "description": "If your software involves blast management, drill data, MWD integration, or exclusion zone management."},
		},
	})

	createSlide(database, s2.ID, "Load & Haul / Mining Operations", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "Load & Haul / Mining Operations",
		"introText": "The biggest department on site. This is where the most people, the most equipment, and the most money are. Excavators load trucks. Trucks haul material. This is where FMS and AHS live.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Run the fleet: excavators, haul trucks, dozers, graders, water carts. Move material from face to dump or crusher. Manage production targets shift by shift."},
			{"title": "Technologies They Use", "description": "Fleet Management Systems (FMS), Autonomous Haulage Systems (AHS), payload monitoring, fatigue detection, collision avoidance systems. This department generates the most operational data on site."},
			{"title": "Key Roles", "description": "Production Supervisor (owns the shift), Excavator Operator, Truck Operator, FMS Dispatcher (control room), AHS Controller."},
			{"title": "When You'd Talk to Them", "description": "If your software touches anything related to fleet activity, production, dispatch, autonomy, or real-time operational data. This is likely your primary stakeholder group."},
		},
	})

	createSlide(database, s2.ID, "Processing & Metallurgy", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Processing & Metallurgy",
		"introText": "Everything downstream of the crusher. They turn ore into saleable product. Usually behind a fence line with its own set of rules.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Crushing, grinding, flotation, leaching, smelting — depends on the commodity. They manage plant throughput, recovery rates, and product quality."},
			{"title": "Technologies They Use", "description": "Process control systems (DCS/SCADA), plant historian databases, metallurgical accounting software. Often completely separate from mine-side IT."},
			{"title": "When You'd Talk to Them", "description": "If your software needs to integrate with plant data (throughput, recovery, grade) or if you're building anything that bridges mine-to-mill."},
		},
	})

	createSlide(database, s2.ID, "Maintenance", "CONTENT", 7, map[string]interface{}{
		"type":      "content",
		"heading":   "Maintenance",
		"introText": "They keep everything running. A single haul truck costs $8-12 million — downtime is measured in thousands of dollars per hour. Maintenance is split into mobile fleet and fixed plant.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Scheduled maintenance (services, component changes), breakdown response, component rebuild, reliability analysis. They manage workshops, parts inventories, and maintenance schedules."},
			{"title": "Technologies They Use", "description": "CMMS (Computerised Maintenance Management Systems — e.g. SAP PM, Pronto, MEX), condition monitoring systems, oil analysis, vibration monitoring, equipment health dashboards."},
			{"title": "When You'd Talk to Them", "description": "If your software touches equipment health, availability, condition monitoring, or maintenance scheduling. The Maintenance Planner is your entry point."},
		},
	})

	createSlide(database, s2.ID, "Safety & Health (HSE)", "CONTENT", 8, map[string]interface{}{
		"type":      "content",
		"heading":   "Safety & Health (HSE)",
		"introText": "They own the rules. Every safety procedure, every incident investigation, every compliance requirement flows through this department. On a mine site, safety outranks production.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Develop and enforce safety management systems, conduct risk assessments, investigate incidents, manage training and inductions, liaise with regulators."},
			{"title": "Technologies They Use", "description": "Safety management systems (INX InControl, SAI360), incident reporting tools, training platforms, fatigue management systems, lone worker tracking."},
			{"title": "When You'd Talk to Them", "description": "If your software involves safety reporting, training delivery, fatigue management, or compliance tracking. They're also relevant if your product is the training platform itself."},
		},
	})

	createSlide(database, s2.ID, "Technology / IT / OT", "CONTENT", 9, map[string]interface{}{
		"type":      "content",
		"heading":   "Technology / IT / OT",
		"introText": "Your primary relationship on site. They own the mine's technology stack — the network, the servers, the onboard systems, and the integrations between them. IT handles business systems; OT (Operational Technology) handles the equipment and control systems.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Manage mine-site networks (WiFi, LTE, UHF), server infrastructure, FMS/AHS platforms, onboard computing, data integration, and cybersecurity. They're the gatekeepers for any software that touches the mine's systems."},
			{"title": "Technologies They Manage", "description": "FMS servers, AHS infrastructure, network switches and radios, onboard truck computers, SCADA systems, data historians, and the integrations between all of them."},
			{"title": "When You'd Talk to Them", "description": "Always. If your software needs to integrate with anything on site, IT/OT controls the handshake. Start here before approaching any other department about data access or system integration."},
		},
	})

	// ── Section 3: People & Technology on Site ──────────────────────────────
	s3 := createSection(database, module.ID, "People & Technology on Site", 2)

	createSlide(database, s3.ID, "People & Technology on Site", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "People & Technology on Site",
		"subtitle": "Who to talk to, where to find them, and how to approach them",
		"bullets": []string{
			"The control room and who sits in it",
			"Key roles that interact with mine technology daily",
			"How to approach each role as a software engineer",
		},
	})

	createSlide(database, s3.ID, "The Control Room", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Control Room",
		"introText": "The nerve centre of the mine. This is where the FMS dispatcher sits, watching every truck on screen in real time. If you want to understand how fleet management works in practice, this is where you start.",
		"items": []map[string]interface{}{
			{"title": "Who's There", "description": "FMS Dispatchers — they assign trucks to excavators, manage queues, handle delays, and keep material moving. On AHS sites, the AHS Controller also sits here or in an adjacent room."},
			{"title": "What They See", "description": "A wall of screens showing every truck's position, status, payload, and assignment. They see the mine in real time — who's loaded, who's hauling, who's queuing, who's broken down."},
			{"title": "How to Approach", "description": "Don't show up during shift changeover — it's the busiest 30 minutes of the day. Ask your host to arrange a quiet time. Observe before you ask questions. The dispatcher is managing a live operation; respect their focus."},
			{"title": "What to Ask", "description": "How they make assignment decisions, what the system does well, what frustrates them, what data they wish they had. These people use the system 12 hours a day — their feedback is gold."},
		},
	})

	createSlide(database, s3.ID, "The Production Supervisor", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "The Production Supervisor",
		"introText": "Owns production for their shift. They're the bridge between the plan (from the mining engineer) and the execution (by operators and dispatchers). They use FMS data constantly to make real-time decisions.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Allocate excavators to dig areas, decide fleet priorities, respond to breakdowns and weather events, and chase tonnage targets. They're accountable for what gets moved on their shift."},
			{"title": "How They Use Technology", "description": "They check FMS dashboards for production rates, truck utilisation, and delays. They use radio and FMS messaging to direct operators. They review shift reports generated from FMS data."},
			{"title": "How to Approach", "description": "Supervisors are busy and time-poor. Be direct, be brief, and have a specific question. Don't ask 'can you tell me about FMS?' — ask 'when you're deciding which trucks to move, what data do you look at?'"},
			{"title": "What to Ask", "description": "What decisions they make using FMS data, what workarounds they've built, where they lose time. Supervisors know where the system fails because they work around it every shift."},
		},
	})

	createSlide(database, s3.ID, "The AHS Controller", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "The AHS Controller",
		"introText": "Manages autonomous trucks from a dedicated control station. AHS doesn't mean 'no humans' — it means different humans, doing different work. The controller is responsible for keeping the autonomous fleet safe and productive.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Monitor autonomous truck behaviour, manage exclusion zones and geofences, handle manual interventions (stopped trucks, obstacle detections), and coordinate with manned equipment entering autonomous areas."},
			{"title": "Mixed Fleet Management", "description": "Most AHS sites run a mix of autonomous and manned trucks. The controller manages the boundary between them — who can go where, how handoffs work, and what happens when a manned vehicle enters an autonomous zone."},
			{"title": "How to Approach", "description": "AHS controllers are highly specialised. If your software touches anything autonomous, this person is your first stop. They understand both the system and its real-world limitations."},
			{"title": "What to Ask", "description": "How they handle edge cases, what causes interventions, how mixed fleet interactions work in practice, what data they need that the system doesn't give them."},
		},
	})

	createSlide(database, s3.ID, "The Mining Engineer", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "The Mining Engineer",
		"introText": "Uses FMS and survey data for planning and reconciliation. They care about the gap between what was planned and what actually happened — and they use technology data to close that gap.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Short-term mine planning, dig plan design, grade control, production reconciliation (planned vs actual), and haul route optimisation."},
			{"title": "How They Use Technology", "description": "They pull data from FMS (tonnes moved, destinations, cycle times), survey systems (GPS pickup, drone surveys), and grade control databases to reconcile and re-plan."},
			{"title": "How to Approach", "description": "Mining engineers are analytical and data-literate. They'll engage deeply if your software solves a data problem they have. Come with specifics, not generalities."},
			{"title": "What to Ask", "description": "Where planned and actual diverge, what data they manually reconcile, what integrations between systems they wish existed. They often build spreadsheet workarounds that reveal unmet needs."},
		},
	})

	createSlide(database, s3.ID, "The Maintenance Planner", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "The Maintenance Planner",
		"introText": "Uses equipment health data to schedule maintenance and maximise fleet availability. If your product touches fleet availability or condition monitoring, start here.",
		"items": []map[string]interface{}{
			{"title": "What They Do", "description": "Schedule services and component changes, manage parts inventory, track equipment hours, and coordinate with operations to get trucks into the workshop without killing production."},
			{"title": "How They Use Technology", "description": "CMMS for work orders and scheduling, condition monitoring systems for oil analysis and vibration data, FMS data for operating hours and utilisation. They often manually bridge data between these systems."},
			{"title": "How to Approach", "description": "Maintenance planners are pragmatic and process-driven. They want solutions that save them time, not more data to look at. Show them how your software reduces manual reconciliation or improves scheduling."},
		},
	})

	createSlide(database, s3.ID, "The IT/OT Team", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "The IT/OT Team",
		"introText": "They own the mine's technology infrastructure. If your software needs to integrate with anything on site, this team controls the handshake. IT handles business systems and networks; OT handles equipment-connected systems and control infrastructure.",
		"items": []map[string]interface{}{
			{"title": "What They Manage", "description": "Site networks (WiFi, private LTE, UHF), FMS and AHS server infrastructure, onboard truck computers, data historians, SCADA systems, and cybersecurity."},
			{"title": "Why They Matter to You", "description": "Every integration goes through them. Data feeds, API access, network configuration, firewall rules, and security approvals — they own all of it. No IT/OT buy-in means no integration."},
			{"title": "How to Approach", "description": "Come with a clear technical ask. They speak your language (software, APIs, networking) but they have hard constraints around security, uptime, and change management. Respect their change control process."},
			{"title": "What to Ask", "description": "What data is available and how to access it, what integration patterns they support, what security and change management processes apply, and what the network architecture looks like."},
		},
	})

	createSlide(database, s3.ID, "What Not to Do", "CONTENT", 7, map[string]interface{}{
		"type":      "content",
		"heading":   "What Not to Do",
		"introText": "Software engineers visiting mines make predictable mistakes. Avoid these and you'll be taken seriously.",
		"items": []map[string]interface{}{
			{"title": "Don't Ask 'Can I See Your FMS?'", "description": "Without context, this is too broad. Know which department your question belongs to. A production question goes to the supervisor, not IT. A data integration question goes to IT/OT, not the control room."},
			{"title": "Don't Show Up During Shift Change", "description": "The 30 minutes around shift changeover are the busiest and most stressful. Supervisors are handing over, dispatchers are re-assigning, and operators are doing pre-starts. Ask your host to avoid this window."},
			{"title": "Don't Assume Miners Are Low-Tech", "description": "Modern mine sites run some of the most complex real-time systems in any industry. Operators, dispatchers, and engineers are highly skilled technology users. Approach them as domain experts, not end users to be educated."},
			{"title": "Don't Skip the Induction", "description": "Even if you think it doesn't apply to you. Even if you've been to another mine. Each site is different and the induction is a legal requirement."},
			{"title": "Don't Treat It Like an Office Visit", "description": "A mine site is a high-consequence environment. The rules aren't suggestions. People have died doing things that seemed harmless. Follow every instruction immediately."},
		},
	})

	createSlide(database, s3.ID, "Shift Handover", "CONTENT", 8, map[string]interface{}{
		"type":      "content",
		"heading":   "Shift Handover",
		"introText": "Shift change is the highest-risk window in a 24-hour mine day. Situational awareness drops, communication loads spike, and people and machines are moving in patterns different to normal production. If you're visiting, know what's happening and stay out of the way.",
		"items": []map[string]interface{}{
			{"title": "Why Shift Change Is High-Risk", "description": "Two crews on site at the same time. Outgoing crew is fatigued. Incoming crew does not yet have current-state awareness. Radio traffic spikes. Accidents disproportionately cluster at shift change — historical data across open-cut mines is consistent on this."},
			{"title": "Pre-Start Meeting", "description": "Before the incoming crew starts equipment, they attend a pre-start briefing with the supervisor. Covers: current ground conditions, active hazards, blast schedules, equipment status, production priorities. Takes 15–30 minutes. No equipment moves until pre-start is complete."},
			{"title": "Equipment Handover", "description": "Incoming operator walks around the machine with the outgoing operator. Covers: running hours, known defects, fuel/oil levels, recent fault codes, any unusual behaviour. Signed off before keys change hands. The walk-around is not optional, even if you 'know' the machine."},
			{"title": "Visitor Rule — Don't Arrive During Shift Change", "description": "Typical shift times are 6am/6pm or 7am/7pm (site-specific). Never arrive at a production area during the ±30 minute window unless your host has briefed the crew. No one has time to host a visitor during handover, and an uninformed visitor is a hazard."},
			{"title": "If You're Caught in Shift Change", "description": "Stay in the crib room or your light vehicle. Do not walk around production areas. Wait for the pre-start to complete and for your host to collect you. It's boring. It's also correct."},
			{"title": "FIFO Context", "description": "On fly-in fly-out sites, shift change also involves personnel movements — buses to and from accommodation, helicopters for remote sites. These movements add traffic, noise, and timing pressure to an already busy period."},
		},
	})

	createSlide(database, s3.ID, "Interacting with Mine Controllers & FMS", "CONTENT", 9, map[string]interface{}{
		"type":      "content",
		"heading":   "Interacting with Mine Controllers & FMS",
		"introText": "The mine controller (sometimes called the dispatcher) sits in the control room and is the real-time operator of the entire mobile fleet. As a software or technology visitor, this is the person you most want to understand — they are the live user of the FMS.",
		"items": []map[string]interface{}{
			{"title": "What a Mine Controller Does", "description": "Manages all mobile equipment movements via the FMS. Assigns trucks to shovels and dumps, manages queues, responds to equipment breakdowns and route closures, coordinates with maintenance and production supervisors. One controller typically manages a fleet of 30–60 trucks plus loading and auxiliary equipment."},
			{"title": "Why They Matter to a Tech Visitor", "description": "The mine controller is your primary window into real-time operations. They use the FMS continuously for 12 hours a shift. They know every quirk, workaround, and pain point. Understanding their day-to-day beats any amount of documentation."},
			{"title": "What They Can Show You", "description": "Live equipment positions on the mine map, current queue lengths at each shovel and crusher, payload history per truck, cycle times and trends, alarm history, recent dispatch decisions. Ask to sit behind them for 30 minutes during a quiet period."},
			{"title": "Protocol", "description": "Always ask permission before entering the control room — it's a restricted area. Don't distract during active incidents (equipment down, blast clearance, emergency response). Keep conversation quiet and out of the controller's line of sight to their screens."},
			{"title": "Common FMS Products", "description": "Modular Mining Dispatch (the original and still most common in open-cut), Wenco (strong in North America and coal), Caterpillar MineStar / Command (tied to Cat fleet), Hexagon Mining. Know which system the site runs before you visit — terminology and workflow differ."},
			{"title": "Good Questions to Ask", "description": "'What's your current bottleneck — shovel, crusher, or dump?' 'How are you managing the shovel queue at Face 7?' 'What would you change about the FMS if you could?' These questions show domain awareness and get real answers."},
		},
	})

	createSlide(database, s3.ID, "Radio Protocol on Site", "CONTENT", 10, map[string]interface{}{
		"type":      "content",
		"heading":   "Radio Protocol on Site",
		"introText": "The mine radio is life-safety communication. Every transmission is heard by everyone on the channel. Protocol is strict because it has to be — a garbled or incorrect call during a critical event can kill someone.",
		"items": []map[string]interface{}{
			{"title": "Treat Every Transmission Seriously", "description": "The radio is not a social channel. No chitchat, no jokes, no unnecessary traffic. Listen before you transmit. Only transmit when you have information the channel needs."},
			{"title": "Call Sign Format", "description": "Your identifier is your vehicle registration or role — 'LV 47' (light vehicle 47), 'Supervisor 3', 'Haul 142', 'Dozer 7'. Never use names. Call sign first, then recipient, then message: 'LV 47 to Controller — requesting entry to Haul Road 7.'"},
			{"title": "Standard Calls", "description": "Location: 'LV 47 on haul road 7 southbound.' Blast clearance: 'All clear 7B — exclusion zone secure.' Hazard: 'Spillage reported haul road 4 at chainage 2100.' Short, specific, unambiguous."},
			{"title": "Road Closures", "description": "Listen for closure announcements before any movement — 'Road 4 closed for blast, reopen 14:30.' Do not enter a closed road. Plan an alternate route and inform the controller. Closures are pushed on the radio, not written down."},
			{"title": "Blast Warnings", "description": "'Fire fire fire — [blast location] — all traffic hold position.' This is a direct, immediate instruction. All vehicles stop where they are. No movement until the all-clear is broadcast. Ignoring this call is grounds for immediate removal from site."},
			{"title": "Visitor Behaviour", "description": "If you're in a light vehicle with a site radio, listen only — do not transmit unless it's an emergency. Your host or driver handles all radio traffic. Even if you hear silence on the channel, someone is listening."},
			{"title": "Emergency Call", "description": "'Mayday Mayday Mayday — [your location] — [nature of emergency].' Used only for imminent danger to life. All other traffic stops. The controller acknowledges and dispatches response. Every person on the channel is expected to stop what they are doing and assist if needed."},
		},
	})

	// ── Section 4: Mine Lingo ───────────────────────────────────────────────
	s4 := createSection(database, module.ID, "Mine Lingo", 3)

	createSlide(database, s4.ID, "Mine Lingo", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Mine Lingo",
		"subtitle": "The language that separates a miner from a software engineer",
		"bullets": []string{
			"Key terms for the pit, the material, and the fleet",
			"Production and FMS/AHS vocabulary",
			"Words that will get you noticed — for the wrong reasons",
		},
	})

	createSlide(database, s4.ID, "The Pit", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "The Pit — Physical Terms",
		"introText": "The open pit is described in precise engineering terms. Using them correctly signals that you've done your homework.",
		"items": []map[string]interface{}{
			{"title": "Bench", "description": "A horizontal step in the pit wall. The mine is excavated in a series of benches, like a giant staircase going down."},
			{"title": "Berm", "description": "A mound of material along the edge of a road or bench, acting as a safety barrier. Berm height must equal the axle height of the largest truck."},
			{"title": "Batter / Batter Angle", "description": "The slope angle of the pit wall between benches. Designed by geotechnical engineers for stability."},
			{"title": "Highwall", "description": "The unexcavated face of the pit — the steep wall you're looking at from the bottom."},
			{"title": "Crest and Toe", "description": "Crest = the top edge of a bench or wall. Toe = the bottom edge. 'Stay back from the crest' means don't go near the edge."},
			{"title": "Ramp / Decline", "description": "The road that connects bench levels, spiralling around the pit. Trucks haul material up the ramp to the surface."},
			{"title": "Pushback / Cutback", "description": "Expanding the pit by mining the walls back. A major planning event that changes the pit geometry."},
		},
	})

	createSlide(database, s4.ID, "The Material", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "The Material — What Gets Moved",
		"introText": "Everything in the pit is classified. Knowing the terminology helps you understand production conversations.",
		"items": []map[string]interface{}{
			{"title": "Ore", "description": "Material with enough valuable mineral to be worth processing. The whole point of the mine."},
			{"title": "Waste", "description": "Material that doesn't contain enough value to process. It still has to be moved — often the majority of material moved on site."},
			{"title": "Overburden", "description": "The material sitting on top of the ore body that must be removed to access it. Often the first thing mined."},
			{"title": "Grade", "description": "The concentration of valuable mineral in the ore. Measured in g/t (grams per tonne) for gold, % for base metals. Higher grade = more valuable."},
			{"title": "Dilution", "description": "Waste material mixed in with ore during mining. Reduces the effective grade. Minimising dilution is a constant operational challenge."},
			{"title": "ROM (Run-of-Mine)", "description": "Ore as it comes out of the pit, before any processing. A ROM pad is a stockpile of unprocessed ore."},
			{"title": "In-Situ", "description": "Material as it exists in the ground, before mining. 'In-situ grade' is the grade before any dilution or loss."},
		},
	})

	createSlide(database, s4.ID, "The Fleet", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "The Fleet — Equipment Names",
		"introText": "Mining equipment has specific names. Using the wrong term is the fastest way to signal you're an outsider.",
		"items": []map[string]interface{}{
			{"title": "Haul Truck", "description": "Not 'truck', not 'vehicle', not 'dumper'. A haul truck. Ultra-class haul trucks carry 200-400 tonnes. They are the backbone of the operation."},
			{"title": "Excavator / Shovel", "description": "The machine that loads trucks. Called an excavator (hydraulic) or a shovel (rope/electric). Never a 'digger' — that's a backhoe in your garden."},
			{"title": "Dozer", "description": "A tracked bulldozer. Used for pushing material, cleaning up around excavators, and maintaining dump faces. 'The dozer is on the face' means it's working at the dig area."},
			{"title": "Grader", "description": "A motor grader — the long machine that smooths road surfaces. Essential for haul road maintenance. 'The grader is on the ramp' means road maintenance is happening."},
			{"title": "Water Cart", "description": "A truck fitted with a water tank and spray bars, used to suppress dust on haul roads. Dust suppression is a safety and environmental requirement."},
			{"title": "Light Vehicle (LV)", "description": "Any vehicle that isn't heavy mobile equipment — Toyota LandCruisers, utes, personnel carriers. 'LV' is the universal term. Never say 'car'."},
			{"title": "Ancillary Equipment", "description": "Support equipment that isn't a haul truck or excavator — dozers, graders, water carts, fuel trucks, tyre handlers. Collectively called 'ancillaries'."},
		},
	})

	createSlide(database, s4.ID, "Production Talk", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Production Talk",
		"introText": "Production conversations use specific metrics. Understanding these helps you follow what supervisors and planners are discussing.",
		"items": []map[string]interface{}{
			{"title": "Cycle Time", "description": "The total time for one truck trip: load + haul + dump + return. The fundamental unit of production measurement."},
			{"title": "Payload", "description": "The weight of material in a truck's tray. Monitored by onboard scales. Under-loading wastes capacity; overloading damages the truck."},
			{"title": "Bucket Factor", "description": "How full the excavator's bucket is relative to its rated capacity. A bucket factor of 0.95 means the bucket is 95% full on average."},
			{"title": "Dig Rate", "description": "How fast the excavator is loading — typically measured in tonnes per operating hour."},
			{"title": "BCMs (Bank Cubic Metres)", "description": "A volume measurement of material in its undisturbed state. Used to measure total material movement. Not the same as loose cubic metres (material expands when broken)."},
			{"title": "Tonnes per Hour", "description": "The rate of material movement. Used for both individual equipment and the mine as a whole."},
		},
	})

	createSlide(database, s4.ID, "FMS and AHS Lingo", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "FMS and AHS Lingo",
		"introText": "These terms come up constantly if you're working in the FMS or AHS space. They're operational, not technical — they describe what happens on the ground.",
		"items": []map[string]interface{}{
			{"title": "Dispatch / Assignment", "description": "The FMS 'dispatches' a truck by assigning it to a destination (an excavator or dump). 'That truck is assigned to EX05' means the FMS told it to go to excavator 5."},
			{"title": "Queuing", "description": "Trucks waiting at the excavator for their turn to load. 'Queuing at the face' means trucks are stacking up — a sign of over-trucking or slow loading."},
			{"title": "Spotting", "description": "The final positioning of a truck next to the excavator before loading. 'Spotting time' is how long it takes to get into position."},
			{"title": "Travel Loaded / Travel Empty", "description": "The two halves of a truck's journey. Loaded = carrying material to the dump. Empty = returning for the next load. FMS tracks both separately."},
			{"title": "Status Codes", "description": "Every truck has a status in the FMS: Operating, Standby, Delay. Each has sub-codes. 'The truck went to delay code 42' — these codes classify what stopped production."},
			{"title": "Intervention (AHS)", "description": "When an autonomous truck stops and requires human input. Interventions are a key AHS performance metric — fewer is better."},
		},
	})

	createSlide(database, s4.ID, "Shift and Time", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Shift and Time",
		"introText": "Mine sites run on shifts, not business hours. The terminology around time and scheduling is specific.",
		"items": []map[string]interface{}{
			{"title": "Day Shift / Night Shift", "description": "Typically 12-hour shifts. Day shift might be 6am-6pm, night shift 6pm-6am. Some sites run 10-hour shifts or 8-hour rotating shifts."},
			{"title": "Pre-Start", "description": "The shift briefing held before work begins. Covers safety topics, production targets, equipment allocations, and any changes from the previous shift. Attendance is mandatory."},
			{"title": "Crib", "description": "A meal or rest break during the shift. 'Crib time' is break time. 'The crib room' is where you eat."},
			{"title": "Changeover", "description": "The handover period between shifts. Information transfer between outgoing and incoming crews. A critical time — don't interrupt."},
			{"title": "Roster", "description": "The work schedule pattern. Common patterns: '2-and-1' (2 weeks on, 1 week off), '8/6' (8 days on, 6 days off), '4/3' (4 on, 3 off). FIFO = Fly In, Fly Out."},
		},
	})

	createSlide(database, s4.ID, "Things That Will Get You Noticed", "CONTENT", 7, map[string]interface{}{
		"type":      "content",
		"heading":   "Things That Will Get You Noticed (Wrong Kind)",
		"introText": "These are the fastest ways to signal you don't belong. Avoid them and you'll be taken more seriously from day one.",
		"items": []map[string]interface{}{
			{"title": "'Truck Driver'", "description": "They're 'operators'. A truck driver delivers groceries. A haul truck operator moves 400 tonnes of rock. The distinction matters to them."},
			{"title": "'Digger'", "description": "It's an excavator or a shovel. A digger is what your nephew uses at the beach."},
			{"title": "'The App'", "description": "On a mine site, software is 'the system' or called by its name (DISPATCH, Jigsaw, MineStar). Calling it 'the app' sounds like you think it's a phone game."},
			{"title": "'Vehicle'", "description": "A haul truck is a haul truck. A light vehicle is an LV. 'Vehicle' is too generic and makes you sound like you're reading a manual."},
			{"title": "'Robot' (for AHS Trucks)", "description": "Autonomous trucks are 'autonomous trucks' or 'AHS trucks'. Calling them robots trivialises the engineering and irritates the people who operate and maintain them."},
			{"title": "Wearing Brand-New PPE", "description": "Not a language mistake, but it marks you as a first-timer. Nothing you can do about it except wear it confidently and not draw attention to it."},
		},
	})

	// ── Section 5: Quiz ─────────────────────────────────────────────────────
	s5 := createSection(database, module.ID, "Knowledge Check", 4)

	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Would You Survive Your First Site Visit?",
		PassMark:         80,
		MaxAttempts:      3,
		ShuffleQuestions: true,
		ShuffleOptions:   true,
		ShowResults:      true,
		SectionID:        &s5.ID,
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
		// ── Multiple Choice (10) ──
		{
			question:    "You're driving a light vehicle on a haul road and see a loaded haul truck approaching. What do you do?",
			explanation: "Loaded haul trucks always have right of way. They weigh up to 600 tonnes and cannot stop quickly. Pull over well clear of the road and wait.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Flash your headlights to signal you're there", false},
				{"Pull over well clear of the road and wait for it to pass", true},
				{"Maintain your lane — you have equal right of way", false},
				{"Speed up to get past before it reaches you", false},
			},
		},
		{
			question:    "Which of the following is NOT mandatory PPE on a typical open pit mine site?",
			explanation: "While gloves may be required in specific areas, the standard mandatory PPE set is hard hat, hi-vis, steel-cap boots, safety glasses, and hearing protection.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Hard hat", false},
				{"Steel-cap boots", false},
				{"Leather gloves", true},
				{"Safety glasses", false},
			},
		},
		{
			question:    "You want to understand how trucks are assigned to excavators in real time. Who do you talk to first?",
			explanation: "The FMS Dispatcher in the control room manages truck assignments in real time. They're the hands-on user of the dispatch system.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"The IT Manager", false},
				{"The FMS Dispatcher in the control room", true},
				{"The Mine General Manager", false},
				{"The Maintenance Planner", false},
			},
		},
		{
			question:    "Which department produces the dig plan that drives load and haul operations?",
			explanation: "Mine Planning / Technical Services designs the pit and produces the short-term and long-term dig plans that drive all downstream operations.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Geology", false},
				{"Processing", false},
				{"Mine Planning / Technical Services", true},
				{"Safety & Health", false},
			},
		},
		{
			question:    "When is the worst time to visit the control room?",
			explanation: "Shift changeover is the busiest and most stressful time in the control room. Dispatchers are re-assigning fleet and handing over to the incoming crew.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"During a blast window", false},
				{"During shift changeover", true},
				{"Mid-morning", false},
				{"After crib", false},
			},
		},
		{
			question:    "Your software needs to pull real-time truck position data from the FMS. Who controls access to that data?",
			explanation: "IT/OT owns the mine's technology infrastructure and controls all data access, API integrations, and network configurations.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"The Production Supervisor", false},
				{"The FMS Dispatcher", false},
				{"The IT/OT Team", true},
				{"The Mining Engineer", false},
			},
		},
		{
			question:    "What does 'queuing at the face' mean?",
			explanation: "Queuing at the face means trucks are waiting in line at the excavator's loading area for their turn to be loaded. It indicates over-trucking or slow loading.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Trucks lining up at the fuel bay", false},
				{"Trucks waiting in line at the excavator to be loaded", true},
				{"Workers queuing for their shift induction", false},
				{"Equipment waiting for maintenance", false},
			},
		},
		{
			question:    "What is a 'pre-start'?",
			explanation: "A pre-start is the mandatory shift briefing held before work begins, covering safety topics, production targets, equipment allocations, and changes from the previous shift.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"A vehicle safety inspection", false},
				{"The shift briefing before work begins", true},
				{"The warm-up period for heavy equipment", false},
				{"The first blast of the day", false},
			},
		},
		{
			question:    "What is the correct term for the slope angle of the pit wall between benches?",
			explanation: "The batter angle is the slope of the pit wall between benches. It's designed by geotechnical engineers to balance ore recovery with wall stability.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Ramp gradient", false},
				{"Batter angle", true},
				{"Bench slope", false},
				{"Crest ratio", false},
			},
		},
		{
			question:    "On a mine site, what does 'LV' stand for?",
			explanation: "LV stands for Light Vehicle — any vehicle that isn't heavy mobile equipment, such as Toyota LandCruisers, utes, and personnel carriers.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Low Voltage", false},
				{"Light Vehicle", true},
				{"Load Volume", false},
				{"Lateral Vision", false},
			},
		},
		// ── True/False (5) ──
		{
			question:    "True or False: On an AHS site, autonomous trucks require no human involvement whatsoever.",
			explanation: "False. AHS trucks require an AHS Controller who monitors behaviour, manages exclusion zones, handles interventions, and coordinates with manned equipment. 'Autonomous' doesn't mean 'unmanaged'.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"True", false},
				{"False", true},
			},
		},
		{
			question:    "True or False: The Mining Engineer is the right person to ask about equipment maintenance scheduling.",
			explanation: "False. Equipment maintenance scheduling is owned by the Maintenance Planner. The Mining Engineer focuses on mine planning, grade control, and production reconciliation.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"True", false},
				{"False", true},
			},
		},
		{
			question:    "True or False: If you hear an emergency call on the UHF radio, you should stop transmitting immediately and listen.",
			explanation: "True. Emergency calls take priority. Stop all non-emergency radio communication and wait for the all-clear.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"True", true},
				{"False", false},
			},
		},
		{
			question:    "True or False: Processing plant control systems are usually managed by the same IT team that manages FMS.",
			explanation: "False. Processing plants often have completely separate control systems (DCS/SCADA) managed by their own technical team. Mine-side IT/OT and plant-side control are frequently independent.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"True", false},
				{"False", true},
			},
		},
		{
			question:    "True or False: 'Ancillary equipment' refers to haul trucks and excavators.",
			explanation: "False. Ancillary equipment refers to support vehicles — dozers, graders, water carts, fuel trucks, tyre handlers. Haul trucks and excavators are primary production equipment.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"True", false},
				{"False", true},
			},
		},
		// ── Scenario-Based (5) ──
		{
			question:    "You arrive on site and realise there are no safety glasses in the visitor PPE pack. Your host says 'it'll be fine, we're only going to the office first.' What should you do?",
			explanation: "Do not proceed to any operational area without full PPE. Even if the first stop is an office, you may need to walk through operational areas. Get the missing PPE before going anywhere.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Proceed — the office area is safe", false},
				{"Ask for safety glasses before proceeding anywhere", true},
				{"Use your sunglasses instead", false},
				{"Go to the pit first to get it over with", false},
			},
		},
		{
			question:    "A supervisor tells you 'we've been losing tonnes to queuing at EX03 all morning.' What is happening?",
			explanation: "Trucks are waiting in line at Excavator 03 for their turn to load. This queuing means trucks sit idle instead of hauling, which directly reduces production tonnes.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Excavator 03 has broken down and is in the maintenance queue", false},
				{"Too many trucks are assigned to Excavator 03 and they're waiting in line to load", true},
				{"The ore grade at Excavator 03 is lower than expected", false},
				{"The haul road to Excavator 03 is blocked", false},
			},
		},
		{
			question:    "You're building software that needs real-time truck data. The Mining Engineer says 'talk to dispatch.' Where on site do you go?",
			explanation: "Dispatch operates from the control room. The FMS Dispatcher manages real-time truck assignments and has the best view of live fleet data. For data access/integration, you'll also need IT/OT.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"The maintenance workshop", false},
				{"The control room", true},
				{"The geology office", false},
				{"The processing plant", false},
			},
		},
		{
			question:    "During a site tour, someone says 'the grader is on the ramp doing road maintenance.' Which piece of equipment and which location are they describing?",
			explanation: "A grader is a motor grader that smooths road surfaces. The ramp is the road connecting bench levels in the pit. Road maintenance on the ramp keeps haul roads safe and reduces tyre wear.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"A bulldozer smoothing the pit floor", false},
				{"A water cart spraying the main access road", false},
				{"A motor grader smoothing the road that connects pit levels", true},
				{"An excavator clearing debris from the bench edge", false},
			},
		},
		{
			question:    "You refer to the autonomous haul trucks as 'the robots' in a meeting with mine operations staff. Why is this a problem?",
			explanation: "Calling AHS trucks 'robots' trivialises the complex engineering and operational skill required to run an autonomous fleet. It signals you don't understand the technology or respect the people who manage it.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"It's technically inaccurate — they're not robots, they're autonomous vehicles", false},
				{"It trivialises the engineering and irritates the people who operate and maintain them", true},
				{"It's a safety violation to use informal terms", false},
				{"It confuses them with underground robotic equipment", false},
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

	// Quiz slide
	createSlide(database, s5.ID, "Site Visit Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Would You Survive Your First Site Visit?",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{"question": "A loaded haul truck is approaching on the haul road. What do you do?", "options": []string{"Flash headlights", "Pull over and wait", "Maintain your lane", "Speed up to pass"}, "correctIndex": 1, "explanation": "Loaded haul trucks always have right of way. Pull over and wait."},
			{"question": "Who do you talk to about real-time truck assignments?", "options": []string{"IT Manager", "FMS Dispatcher", "Mine GM", "Maintenance Planner"}, "correctIndex": 1, "explanation": "The FMS Dispatcher in the control room manages truck assignments in real time."},
			{"question": "What does 'LV' stand for?", "options": []string{"Low Voltage", "Light Vehicle", "Load Volume", "Lateral Vision"}, "correctIndex": 1, "explanation": "LV = Light Vehicle — any non-heavy vehicle on site."},
		},
	})

	// Completion slide
	createSlide(database, s5.ID, "Module Complete", "COMPLETION", 1, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "You're ready for your first mine site visit. Remember: respect the rules, know the people, speak the language.",
		"keyTakeaways": []string{
			"PPE is non-negotiable — hard hat, hi-vis, steel caps, safety glasses, hearing protection",
			"Haul trucks always have right of way — you are invisible to them",
			"Know the value chain: Geology → Planning → Drill & Blast → Load & Haul → Processing",
			"The control room is the nerve centre — FMS dispatchers and AHS controllers sit here",
			"IT/OT controls all technology integration — start there for data access",
			"Use the right terminology — operators, excavators, haul trucks, LVs",
		},
		"showCertificate": true,
	})

	log.Println("Visiting a Mine Site module seeded with 5 sections, 30+ slides, and 20-question quiz")
}
