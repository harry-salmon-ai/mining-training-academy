package seed

import (
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDrillAndBlastModule(database *gorm.DB, force bool) {
	var moduleCount int64
	database.Model(&models.Module{}).Where("slug = ?", "drill-and-blast-operations").Count(&moduleCount)
	if moduleCount > 0 {
		if !force {
			log.Println("Drill & Blast Operations module already exists, skipping")
			return
		}
		forceDeleteModule(database, "drill-and-blast-operations")
	}

	log.Println("Seeding Drill & Blast Operations module...")

	var category models.Category
	database.Where("slug = ?", "drill-blast").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed Drill & Blast Operations: category or admin not found")
		return
	}

	desc := "A deep dive into drill and blast operations in open-cut mining — covering drilling equipment and pattern design, explosive types and loading techniques, blast design and timing, execution safety, post-blast fragmentation assessment, and the mine-to-mill concept that links blast quality to total value chain performance."
	dur := 90
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Drill & Blast Operations",
		Slug:        "drill-and-blast-operations",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelFoundation,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"drill and blast", "blast design", "powder factor", "ANFO", "emulsion", "electronic detonators", "fragmentation", "flyrock", "mine to mill", "blast pattern", "stemming", "burden", "spacing", "DTH", "rotary drill"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 1: Introduction ─────────────────────────────────────────
	s1 := createSection(database, module.ID, "Introduction", 0)

	createSlide(database, s1.ID, "Drill & Blast Operations", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Drill & Blast Operations",
		"subtitle": "Fragmenting rock to feed the mining value chain",
		"bullets": []string{
			"How drill and blast quality determines everything downstream — load, haul, crush, process",
			"Drilling equipment, pattern design, and performance measurement",
			"Explosive types, blast design, and timing sequences",
			"Safety, execution, post-blast assessment, and the mine-to-mill concept",
		},
	})

	createSlide(database, s1.ID, "Drill & Blast in the Mining Value Chain", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Drill & Blast in the Mining Value Chain",
		"introText": "Drill and blast is the first physical act of mining. The quality of fragmentation at this step cascades through every subsequent process — loading productivity, haul truck payload, crusher throughput, and mill energy consumption. A well-designed blast can reduce total mining cost by 10–15%; a poor blast inflates costs at every downstream step and may not be recoverable.",
		"items": []map[string]interface{}{
			{
				"title":       "Loading Productivity",
				"description": "Well-fragmented rock fills excavator buckets in fewer passes. Oversized material causes bucket stalling, reduced fill factors (0.7 vs 0.95), and longer load times. A 10% improvement in fragmentation can deliver a 15–20% improvement in excavator productivity.",
			},
			{
				"title":       "Haul Truck Payload",
				"description": "Consistent fragmentation means consistent bucket fills, which means consistent truck payloads. Oversized fragments create voids in the truck body, reducing payload by 5–15% per load.",
			},
			{
				"title":       "Crusher and Mill Throughput",
				"description": "Oversized material causes crusher bridging, increased energy consumption, and reduced throughput. The primary crusher is typically the bottleneck — feeding it well-fragmented material is the highest-value outcome of a blast.",
			},
			{
				"title":       "The Cost Multiplier",
				"description": "Drill and blast represents 15–20% of total mining cost, but it determines the cost of the 80–85% that follows. Every dollar under-spent on blast quality costs $3–5 in downstream inefficiency.",
			},
		},
	})

	createSlide(database, s1.ID, "The Drill & Blast Process", "PROCESS", 2, map[string]interface{}{
		"type":        "process",
		"heading":     "The Drill & Blast Process",
		"description": "The end-to-end sequence from bench preparation through to post-blast assessment. Each step creates the conditions for the next — shortcutting any step degrades the final result.",
		"steps": []map[string]interface{}{
			{
				"title":       "Bench Preparation",
				"description": "Survey the bench, mark out the drill pattern, confirm grade control boundaries. The bench must be clean and level for accurate collar positioning.",
			},
			{
				"title":       "Drilling",
				"description": "Drill holes to design depth at the specified pattern geometry (burden × spacing). Record penetration rates and ground conditions.",
			},
			{
				"title":       "Charging (Loading)",
				"description": "Load explosives into each hole per the blast design. Explosive type, column charge, deck charging, and stemming height are all specified.",
			},
			{
				"title":       "Timing & Hookup",
				"description": "Connect detonators in the designed sequence. Timing controls rock movement direction, fragmentation, and vibration.",
			},
			{
				"title":       "Firing",
				"description": "Clear the blast area, confirm exclusion zone, fire. Monitor for misfires.",
			},
			{
				"title":       "Post-Blast Assessment",
				"description": "Survey the muckpile, assess fragmentation, measure heave, check for misfires, and feed results back to the next blast design.",
			},
		},
	})

	createSlide(database, s1.ID, "Where Drill & Blast Sits", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Where Drill & Blast Sits",
		"description": "Drill and blast sits between mine planning and load & haul. It converts the mine plan's geometric design into physical broken rock that the fleet can move. The quality of this conversion determines the productivity of everything that follows.",
		"diagramType": "process",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"label": "Mine Planning", "description": "Bench design, grade control, scheduling", "highlight": false},
				{"label": "Drill & Blast", "description": "Fragment rock to target size", "highlight": true},
				{"label": "Load & Haul", "description": "Excavator loads, trucks transport", "highlight": false},
				{"label": "Crush & Process", "description": "Primary/secondary crushing, milling", "highlight": false},
				{"label": "Product", "description": "Concentrate, metal, coal to market", "highlight": false},
			},
			"callout": "Drill & Blast = 15–20% of cost, but controls 80–85% of downstream efficiency",
		},
	})

	// ── Section 2: Drilling Equipment & Methods ─────────────────────────
	s2 := createSection(database, module.ID, "Drilling Equipment & Methods", 1)

	createSlide(database, s2.ID, "Blast Hole Drill Rig", "EQUIPMENT", 0, map[string]interface{}{
		"type":        "equipment",
		"heading":     "Rotary Blast Hole Drill",
		"description": "The rotary blast hole drill is the workhorse of open-cut blast hole drilling. It drills vertical or near-vertical holes using a tri-cone roller bit and rotary action with compressed air for cuttings removal.",
		"model":       "Rotary Blast Hole Drill",
		"specs": []map[string]interface{}{
			{"label": "Hole Diameter", "value": "171–311 mm (6.75–12.25 in)"},
			{"label": "Hole Depth", "value": "8–20 m (single pass)"},
			{"label": "Penetration Rate", "value": "15–45 m/hr (rock dependent)"},
			{"label": "Operating Weight", "value": "50–120 tonnes"},
			{"label": "Compressor Capacity", "value": "25–50 m³/min"},
			{"label": "Bit Life", "value": "300–2,000 m (varies with abrasivity)"},
			{"label": "Engine Power", "value": "500–1,000 kW"},
			{"label": "Mast Height", "value": "12–20 m"},
		},
		"advantages": []string{
			"High penetration rates in soft to medium rock",
			"Large hole diameters for high powder factor designs",
			"Excellent for production drilling in open-cut benches",
			"MWD (Measure While Drilling) integration for rock characterisation",
		},
		"limitations": []string{
			"Reduced penetration in very hard, competent rock (UCS > 200 MPa)",
			"Requires level, stable bench surface",
			"High capital cost ($3–8M per unit)",
			"Limited hole angle flexibility (typically 0–20° from vertical)",
		},
	})

	createSlide(database, s2.ID, "Rotary vs Down-the-Hole (DTH) Drilling", "COMPARISON", 1, map[string]interface{}{
		"type":        "comparison",
		"heading":     "Rotary vs Down-the-Hole (DTH) Drilling",
		"description": "The two primary drilling methods in open-cut mining. Selection depends on rock hardness, hole diameter requirements, and production volume.",
		"categories": []map[string]interface{}{
			{
				"name":        "Rotary Drilling",
				"description": "Tri-cone roller bit rotated by the drill string. Compressed air flushes cuttings. Best for soft to medium rock in high-production environments.",
				"pros": []string{
					"Highest penetration rates in soft-medium rock (20–45 m/hr)",
					"Large hole diameters (up to 311 mm)",
					"Continuous drilling (no hammer cycle)",
					"Excellent for consistent bench geology",
				},
				"cons": []string{
					"Penetration drops sharply in hard rock (UCS > 200 MPa)",
					"Bit wear is high in abrasive formations",
					"Hole deviation increases with depth in fractured ground",
					"Higher fuel/energy consumption",
				},
				"specs": map[string]interface{}{
					"Hole Diameter": "171–311 mm",
					"Penetration":   "15–45 m/hr",
					"Best Rock UCS":  "50–200 MPa",
				},
			},
			{
				"name":        "Down-the-Hole (DTH) Drilling",
				"description": "Pneumatic percussion hammer at the bottom of the hole delivers impact energy directly to the rock face. Energy transfer efficiency does not decrease with depth.",
				"pros": []string{
					"Consistent penetration regardless of hole depth",
					"Superior performance in hard rock (UCS > 150 MPa)",
					"Straighter holes with less deviation",
					"Lower energy cost per metre in hard formations",
				},
				"cons": []string{
					"Lower penetration rates than rotary in soft rock",
					"Limited hole diameter range (typically 89–254 mm)",
					"Higher bit cost per metre in abrasive rock",
					"Compressor capacity is the production bottleneck",
				},
				"specs": map[string]interface{}{
					"Hole Diameter": "89–254 mm",
					"Penetration":   "8–25 m/hr",
					"Best Rock UCS":  "150–350 MPa",
				},
			},
		},
	})

	createSlide(database, s2.ID, "Drill Pattern Design", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Drill Pattern Design",
		"introText": "The drill pattern determines the geometric distribution of explosive energy in the rock mass. Every parameter — burden, spacing, sub-drill, stemming — has a design value derived from rock properties, explosive type, and the desired fragmentation outcome. Getting the pattern right is the single highest-leverage decision in blast design.",
		"items": []map[string]interface{}{
			{
				"title":       "Burden (B)",
				"description": "The distance from the blast hole to the nearest free face. Too small: excessive flyrock and airblast. Too large: poor fragmentation, toe problems, and incomplete breakage at the face. Typical range: 3.5–8.0 m depending on hole diameter and explosive type. Rule of thumb: B = 25–40 × hole diameter.",
			},
			{
				"title":       "Spacing (S)",
				"description": "The distance between holes in a row. Controls lateral energy distribution. Typical S/B ratio: 1.0–1.3 for staggered patterns, 1.15 for equilateral triangle patterns. Under-spaced = excessive crushing near holes and wasted energy. Over-spaced = coarse fragmentation between holes.",
			},
			{
				"title":       "Sub-drill",
				"description": "Hole depth below the design bench floor. Required to ensure clean breakage at bench grade. Typical: 0.3 × burden (0.5–2.0 m). Excessive sub-drill wastes explosives and can damage the bench below.",
			},
			{
				"title":       "Stemming",
				"description": "Inert material (drill cuttings or crushed aggregate) placed in the top of the hole above the explosive column. Confines explosive gases to improve energy transfer. Typical: 0.7–1.0 × burden. Insufficient stemming causes airblast and flyrock from the collar.",
			},
			{
				"title":       "Bench Height",
				"description": "The vertical distance between bench surfaces. Determines hole depth (bench height + sub-drill). Typical open-cut: 10–15 m. Higher benches increase powder efficiency but require larger equipment and carry greater face failure risk.",
			},
		},
	})

	createSlide(database, s2.ID, "Blast Pattern Layout", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Blast Pattern Layout — Burden, Spacing, and Echelon Timing",
		"description": "Top-down view of a production blast pattern. Holes are laid out on a staggered grid defined by burden (B — distance to the free face) and spacing (S — distance between holes in a row). The echelon timing sequence — shown as contour lines of equal firing time — creates a moving free face: each hole fires with enough delay after its neighbours for the broken rock to start moving, giving the next hole somewhere to throw material. This combination of geometry and timing is what turns a static pattern into a controlled fragmentation event.",
		"diagramType": "image",
		"imageUrl":    "/assets/diagrams/blast_pattern_top_down.png",
	})

	// ── Section 3: Explosives and Loading ────────────────────────────────
	s3 := createSection(database, module.ID, "Explosives and Loading", 2)

	createSlide(database, s3.ID, "Explosive Types in Open-Cut Mining", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Explosive Types in Open-Cut Mining",
		"introText": "Explosive selection is driven by rock properties, water conditions, hole diameter, and required energy output. The three dominant product families in open-cut mining are ANFO, emulsion, and heavy ANFO blends. Each has different energy density, water resistance, and detonation characteristics.",
		"items": []map[string]interface{}{
			{
				"title":       "ANFO (Ammonium Nitrate / Fuel Oil)",
				"description": "The simplest and cheapest bulk explosive. A 94:6 mix of porous ammonium nitrate prills and diesel fuel. Energy density: ~2.5 MJ/kg. Density: 0.8–0.85 g/cc. Zero water resistance — cannot be used in wet holes. Accounts for ~70% of all commercial explosives used in mining by volume. Cost: ~$0.30–0.50/kg delivered.",
			},
			{
				"title":       "Emulsion",
				"description": "A water-in-oil emulsion of ammonium nitrate solution with a fuel phase. Inherently water resistant. Energy density: ~3.0–3.5 MJ/kg. Density: 1.0–1.25 g/cc (pumpable). Used in wet holes, small-diameter holes, and where higher energy density is needed. More expensive than ANFO (~$0.80–1.50/kg).",
			},
			{
				"title":       "Heavy ANFO (HA)",
				"description": "A blend of ANFO and emulsion, typically 30–70% emulsion by weight. Combines the low cost of ANFO with improved water resistance and higher energy density. Density: 0.9–1.3 g/cc (adjustable). The dominant product in large open-cut operations where some water is present but full emulsion is not economically justified.",
			},
			{
				"title":       "Sensitised Emulsion (Bulk)",
				"description": "Emulsion that has been chemically sensitised (gassed) to create micro-bubbles within the emulsion matrix. This is the pumpable product delivered by mobile manufacturing units (MMUs) directly to the hole. Density and energy can be precisely controlled by adjusting the gassing agent concentration.",
			},
		},
	})

	createSlide(database, s3.ID, "ANFO vs Emulsion — When to Use Each", "COMPARISON", 1, map[string]interface{}{
		"type":        "comparison",
		"heading":     "ANFO vs Emulsion — When to Use Each",
		"description": "The choice between ANFO and emulsion is not about which is 'better' — it is about matching the explosive to the conditions. Using ANFO in wet holes wastes money and creates misfires. Using full emulsion in dry, soft rock wastes money without improving results.",
		"categories": []map[string]interface{}{
			{
				"name":        "ANFO",
				"description": "Simple, low-cost bulk explosive for dry-hole production blasting.",
				"pros": []string{
					"Lowest cost per kg of explosive energy",
					"Simple to manufacture and deliver",
					"Excellent performance in dry conditions",
					"Easy to control loading density",
				},
				"cons": []string{
					"Zero water resistance — desensitises and fails in standing water",
					"Lower energy density than emulsion",
					"Creates more NOx fumes than emulsion",
					"Cannot be pumped (auger-loaded only)",
				},
				"specs": map[string]interface{}{
					"Energy":  "2.5 MJ/kg",
					"Density": "0.8 g/cc",
					"Cost":    "$0.30–0.50/kg",
					"VOD":     "3,500–4,500 m/s",
				},
			},
			{
				"name":        "Emulsion",
				"description": "Water-resistant pumpable explosive for wet or high-energy applications.",
				"pros": []string{
					"Fully water resistant — reliable in standing water",
					"Higher energy density (up to 40% more than ANFO)",
					"Pumpable delivery via MMU for fast loading",
					"Lower fume classification (better for confined areas)",
				},
				"cons": []string{
					"2–3× the cost of ANFO per kg",
					"Requires mobile manufacturing unit (MMU) for bulk delivery",
					"Over-confinement in hard rock can cause excessive vibration",
					"Shelf life limitations — sensitised product must be used within hours",
				},
				"specs": map[string]interface{}{
					"Energy":  "3.0–3.5 MJ/kg",
					"Density": "1.0–1.25 g/cc",
					"Cost":    "$0.80–1.50/kg",
					"VOD":     "4,500–6,000 m/s",
				},
			},
		},
	})

	createSlide(database, s3.ID, "Detonator Systems", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Detonator Systems",
		"introText": "Detonators control when each hole fires. The timing sequence determines fragmentation quality, muckpile shape, vibration levels, and whether the blast moves rock in the intended direction. Modern electronic detonators have transformed blast design from a constrained timing grid to a fully programmable firing sequence.",
		"items": []map[string]interface{}{
			{
				"title":       "Non-Electric (Nonel) Detonators",
				"description": "Signal transmitted via a shockwave in a hollow plastic tube. Timing set by pre-manufactured delay elements (25 ms, 42 ms, 65 ms, etc.). Reliable, simple, and widely used. Limitation: fixed delay increments and timing scatter of ±5–10 ms, which limits pattern complexity.",
			},
			{
				"title":       "Electronic Detonators",
				"description": "Programmable digital delay with 1 ms precision and near-zero scatter (<0.5 ms). Each detonator is individually addressable and can be programmed in the field via a handheld unit. Enables complex timing designs previously impossible: curved firing fronts, variable inter-hole delays, and fine-tuned vibration management. Cost: 3–5× nonel on a per-unit basis.",
			},
			{
				"title":       "Why Timing Precision Matters",
				"description": "A 5 ms scatter in a 25 ms inter-hole delay means adjacent holes can fire almost simultaneously (20 ms gap) or with excessive delay (30 ms gap). The result is inconsistent fragmentation, increased vibration, and unpredictable muckpile geometry. Electronic detonators remove this variability.",
			},
			{
				"title":       "Detonator Selection Decision",
				"description": "Nonel is appropriate for routine production blasts in homogeneous rock with no vibration constraints. Electronic is required when: vibration limits are tight (near infrastructure or communities), complex timing is needed, fragmentation consistency is critical to downstream performance, or when the cost of a poor blast exceeds the detonator cost premium.",
			},
		},
	})

	createSlide(database, s3.ID, "Charge Loading and Stemming", "PROCESS", 3, map[string]interface{}{
		"type":        "process",
		"heading":     "Charge Loading and Stemming",
		"description": "The loading process determines how explosive energy is distributed within each hole. Every step from priming to stemming affects blast performance.",
		"steps": []map[string]interface{}{
			{
				"title":       "Verify Hole Condition",
				"description": "Measure hole depth (compare to design). Check for water (if wet, determine water column height). Inspect for caving or obstructions. A hole that has collapsed or is shorter than design must be re-drilled or the charge design adjusted.",
			},
			{
				"title":       "Place the Primer",
				"description": "Lower the primer cartridge (booster + detonator assembly) to the bottom of the hole (bottom priming) or the top of the charge column (top priming). Bottom priming is standard in open-cut — it maximises confinement and ensures the entire column detonates upward.",
			},
			{
				"title":       "Load the Explosive Column",
				"description": "Pour or pump explosive to the designed column height. For ANFO: auger delivery from a bulk truck. For emulsion/heavy ANFO: pump delivery from an MMU. Column height = hole depth minus stemming height minus sub-drill (which is already drilled below grade).",
			},
			{
				"title":       "Deck Charging (If Specified)",
				"description": "In some designs, an inert gap (deck) separates two explosive columns within the same hole. Deck charging reduces peak vibration by splitting the charge into two separately timed events. Used near vibration-sensitive receivers.",
			},
			{
				"title":       "Stem the Hole",
				"description": "Fill the remaining hole length above the explosive column with inert stemming material (clean crushed aggregate, 10–20 mm). Stemming confines explosive gases, improving energy transfer to the rock. Drill cuttings are poor stemming material — angular aggregate interlocks better and provides superior confinement.",
			},
			{
				"title":       "Record and Verify",
				"description": "Record the as-loaded charge weight, stemming height, and any variations from the design. Every hole must be reconciled against the blast design before hookup begins.",
			},
		},
	})

	createSlide(database, s3.ID, "MMU Loading a Blast Hole", "IMAGE", 4, map[string]interface{}{
		"type":        "image",
		"heading":     "Mobile Manufacturing Unit (MMU) Loading a Blast Hole",
		"description": "A mobile manufacturing unit pumps sensitised emulsion directly into a blast hole on the bench. The MMU carries separate raw materials (ammonium nitrate solution, fuel phase, and gassing agent) and mixes the final product at the hole — this is safer than transporting pre-sensitised explosives. Loading rates of 400–800 kg/minute allow a production blast of 100+ holes to be loaded in a single shift.",
		"imageUrl":    "https://assets.miningacademy.com/images/mmu-loading-blast-hole.jpg",
	})

	// ── Section 4: Blast Design ─────────────────────────────────────────
	s4 := createSection(database, module.ID, "Blast Design", 3)

	createSlide(database, s4.ID, "Powder Factor and Energy Distribution", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Powder Factor and Energy Distribution",
		"introText": "Powder factor is the ratio of explosive mass to rock mass (kg/t) or explosive mass to rock volume (kg/m³). It is the primary design parameter that links explosive energy to the desired fragmentation outcome. But powder factor alone is insufficient — how that energy is distributed through the pattern geometry matters just as much.",
		"items": []map[string]interface{}{
			{
				"title":       "Powder Factor Definition",
				"description": "Powder factor (PF) = total explosive mass ÷ total rock volume (kg/m³). Typical open-cut PF: 0.3–0.8 kg/m³ (or 0.12–0.32 kg/t for rock density 2.5 t/m³). Higher PF = finer fragmentation but higher explosive cost. The economic optimum is the PF that minimises total mining cost, not blast cost.",
			},
			{
				"title":       "Energy Factor vs Powder Factor",
				"description": "Different explosives deliver different energy per kg. Comparing ANFO at 2.5 MJ/kg to emulsion at 3.5 MJ/kg: a 0.5 kg/m³ ANFO blast delivers 1.25 MJ/m³, while a 0.5 kg/m³ emulsion blast delivers 1.75 MJ/m³. Energy factor (MJ/m³) is a more accurate design parameter than powder factor when comparing different explosive products.",
			},
			{
				"title":       "Distribution Matters",
				"description": "The same total explosive weight distributed across more, smaller holes fragments rock more uniformly than fewer, larger holes. This is why pattern design (burden × spacing) and powder factor are inseparable — you cannot set one without the other.",
			},
			{
				"title":       "The Total Cost Curve",
				"description": "Increasing powder factor increases blast cost but decreases downstream costs (loading, hauling, crushing). The minimum total cost is typically at a higher PF than the minimum blast cost. This is the economic argument for spending more on drill and blast to save more in load, haul, and processing.",
			},
		},
	})

	createSlide(database, s4.ID, "Powder Factor Calculator", "INTERACTIVE", 1, map[string]interface{}{
		"type":            "interactive",
		"heading":         "Blast Powder Factor Calculator",
		"description":     "Adjust pattern geometry and charge parameters to see how they affect powder factor and explosive consumption per blast.",
		"interactionType": "calculator",
		"controls": []map[string]interface{}{
			{"id": "burden", "label": "Burden (m)", "type": "range", "min": 3.0, "max": 8.0, "step": 0.5, "default": 5.0},
			{"id": "spacing", "label": "Spacing (m)", "type": "range", "min": 3.5, "max": 10.0, "step": 0.5, "default": 6.0},
			{"id": "benchHeight", "label": "Bench Height (m)", "type": "range", "min": 8, "max": 18, "step": 1, "default": 12},
			{"id": "holeDensity", "label": "Explosive Density (g/cc)", "type": "range", "min": 0.8, "max": 1.3, "step": 0.05, "default": 0.85},
			{"id": "holeDiameter", "label": "Hole Diameter (mm)", "type": "range", "min": 127, "max": 311, "step": 10, "default": 200},
			{"id": "stemmingHeight", "label": "Stemming Height (m)", "type": "range", "min": 2.0, "max": 5.0, "step": 0.5, "default": 3.5},
		},
		"outputs": []map[string]interface{}{
			{"id": "chargeLength", "label": "Charge Length", "formula": "benchHeight - stemmingHeight + 1.5", "unit": "m", "description": "Column length including sub-drill (~1.5 m)"},
			{"id": "chargeWeight", "label": "Charge per Hole", "formula": "holeDensity * 3.14159 * (holeDiameter/2000) * (holeDiameter/2000) * (benchHeight - stemmingHeight + 1.5) * 1000", "unit": "kg"},
			{"id": "rockVolume", "label": "Rock per Hole", "formula": "burden * spacing * benchHeight", "unit": "m³"},
			{"id": "powderFactor", "label": "Powder Factor", "formula": "(holeDensity * 3.14159 * (holeDiameter/2000) * (holeDiameter/2000) * (benchHeight - stemmingHeight + 1.5) * 1000) / (burden * spacing * benchHeight)", "unit": "kg/m³"},
		},
		"feedback": []map[string]interface{}{
			{"condition": "powderFactor < 0.3", "message": "Powder factor below 0.3 kg/m³ — expect coarse fragmentation, poor digging, and potential toe problems. Consider tighter pattern or higher energy explosive.", "type": "warning"},
			{"condition": "powderFactor > 0.8", "message": "Powder factor above 0.8 kg/m³ — excessive energy. Risk of flyrock, airblast, and unnecessary explosive cost. Review pattern geometry.", "type": "warning"},
			{"condition": "powderFactor >= 0.3 && powderFactor <= 0.8", "message": "Powder factor in typical open-cut range. Fine-tune based on fragmentation results and downstream performance.", "type": "success"},
		},
	})

	createSlide(database, s4.ID, "Timing and Sequencing", "DIAGRAM", 2, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Blast Timing Sequence — Why Order Matters",
		"description": "The firing sequence controls which holes fire first, how the rock moves, and what free face is available at each moment of the blast. Without a free face, explosive energy has nowhere to move the rock — it converts to vibration and airblast instead of fragmentation.",
		"diagramType": "process",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"label": "Free Face Established", "description": "The first row fires toward the open bench face. This row has a free face — rock can move horizontally into the void."},
				{"label": "Sequential Row Firing", "description": "Each subsequent row fires after a designed inter-row delay (typically 25–65 ms per row). By the time each row fires, the previous row has begun moving forward, creating a new free face."},
				{"label": "Inter-Hole Timing", "description": "Within each row, holes fire with a designed inter-hole delay (typically 5–25 ms). This controls fragmentation quality and limits instantaneous charge (vibration)."},
				{"label": "Rock Movement Direction", "description": "The timing sequence directs rock movement. Firing from left to right moves the muckpile to the right. V-patterns create a compact muckpile. Echelon patterns direct material toward the excavator."},
				{"label": "Vibration Management", "description": "Maximum instantaneous charge (MIC) = the largest charge weight detonating within any 8 ms window. Timing controls MIC. Electronic detonators allow MIC to be precisely designed, not estimated."},
			},
		},
	})

	createSlide(database, s4.ID, "Blast Design Software and Modelling", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Blast Design Software and Modelling",
		"introText": "Modern blast design is computer-aided. Software tools model the interaction between pattern geometry, explosive energy, rock properties, and timing to predict fragmentation, muckpile shape, vibration, and flyrock. The shift from experience-based to data-driven blast design has fundamentally improved consistency.",
		"items": []map[string]interface{}{
			{
				"title":       "Blast Design Software",
				"description": "Tools like Orica SHOTPlus, BME BlastMap, AEGIS, and Dyno Nobel DynoConsult provide 3D blast modelling with timing simulation, fragmentation prediction, vibration prediction, and flyrock risk assessment. The software takes pattern geometry, explosive properties, and rock parameters as inputs.",
			},
			{
				"title":       "Rock Mass Characterisation",
				"description": "Blast design quality depends on knowing what you are blasting. Rock mass properties include UCS (unconfined compressive strength), fracture frequency, joint orientation, density, and abrasivity. MWD data from drilling provides real-time rock property information that can adjust the blast design hole-by-hole.",
			},
			{
				"title":       "Fragmentation Prediction Models",
				"description": "The Kuz-Ram model is the most widely used empirical fragmentation model. Inputs: rock factor, powder factor, hole diameter, spacing, burden. Output: predicted size distribution curve. Useful for comparative design but underpredicts fines. The Swebrec function provides better fit for the fine end of the distribution.",
			},
			{
				"title":       "Iterative Design Improvement",
				"description": "The best blast designs are iterative: design, fire, measure fragmentation, compare to prediction, adjust parameters, repeat. Without post-blast measurement, there is no feedback loop — and without feedback, blast quality cannot systematically improve.",
			},
		},
	})

	// ── Section 5: Blast Execution and Safety ────────────────────────────
	s5 := createSection(database, module.ID, "Blast Execution and Safety", 4)

	createSlide(database, s5.ID, "Blast Clearance Procedure", "PROCESS", 0, map[string]interface{}{
		"type":        "process",
		"heading":     "Blast Clearance Procedure",
		"description": "Blast clearance is the single most safety-critical procedure in drill and blast operations. Every person and every machine must be accounted for and confirmed at a safe distance before firing. Clearance failures are the primary cause of blast-related fatalities.",
		"steps": []map[string]interface{}{
			{
				"title":       "Blast Area Defined",
				"description": "The blast engineer defines the exclusion zone based on flyrock risk assessment. Typical minimum: 500 m from the blast perimeter for personnel, 300 m for light vehicles behind solid cover. Flyrock exclusion zones can extend to 800+ m for large blasts or adverse geometry.",
			},
			{
				"title":       "Pre-Blast Communication",
				"description": "Blast notification issued to all site personnel via radio, FMS alert, and siren system. Notification includes blast time, exclusion zone boundaries, and expected all-clear time. Minimum 30-minute warning for planned blasts.",
			},
			{
				"title":       "Sentry Placement",
				"description": "Physical sentries posted at all access points to the exclusion zone. Sentries must be in position and confirmed by the shotfirer before the countdown begins. No person may enter the exclusion zone after sentries are posted.",
			},
			{
				"title":       "Personnel and Equipment Count",
				"description": "A physical count of all personnel and equipment inside the exclusion zone. The FMS vehicle tracking system provides a digital count, but the shotfirer must confirm by radio that every person is accounted for. The count must reconcile to zero inside the zone.",
			},
			{
				"title":       "Final Siren Sequence",
				"description": "Standard Australian protocol: 2-minute warning (continuous siren), 1-minute warning (intermittent siren), firing (long blast). Some sites use 5-2-1 minute warning sequences. The sequence is non-negotiable and must be completed even if the blast is delayed.",
			},
			{
				"title":       "Fire Command",
				"description": "The shotfirer (licensed blaster) confirms all sentries in position, all personnel clear, exclusion zone secure. Only the shotfirer can initiate the blast. A misfire protocol must be in place before the fire command is given.",
			},
			{
				"title":       "All-Clear",
				"description": "After firing, the shotfirer inspects the blast site for misfires, unexploded material, and flyrock before issuing the all-clear. Minimum 5-minute post-blast wait (some jurisdictions require longer). No person enters the blast area until the all-clear is broadcast.",
			},
		},
	})

	createSlide(database, s5.ID, "Flyrock and Airblast Hazards", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Flyrock and Airblast Hazards",
		"introText": "Flyrock and airblast are the two mechanisms by which a blast can cause harm beyond the blast pattern itself. Flyrock kills; airblast damages infrastructure and can cause injury. Both are preventable through design, but their prevention requires understanding the mechanisms that cause them.",
		"items": []map[string]interface{}{
			{
				"title":       "Flyrock Mechanisms",
				"description": "Flyrock is rock propelled by blast gases through three mechanisms: (1) cratering — insufficient stemming allows gas to vent through the collar, ejecting stemming and rock; (2) face burst — excessive burden or inadequate free face causes gas to break through the face, launching material; (3) rifling — blast gases follow weak zones (joints, previous blast damage) and eject material at unexpected trajectories. Flyrock distances of 500–1,000 m are documented.",
			},
			{
				"title":       "Flyrock Prevention",
				"description": "Correct stemming height (minimum 0.7 × burden with angular aggregate), accurate burden control (survey the face before designing), sub-drill within design limits, and crest control (no holes drilled within 2 m of the crest without mitigation). Face profiling via drone survey enables precise burden measurement hole-by-hole.",
			},
			{
				"title":       "Airblast",
				"description": "The air pressure wave generated by a blast. Sources: gas venting through stemming (the dominant source), ground-coupled air wave, and direct rock displacement. Measured in dB(L). Regulatory limit typically 115–120 dB(L) at nearest residence. Airblast attenuates predictably with distance — the design tool is stemming quality and confinement.",
			},
			{
				"title":       "Vibration",
				"description": "Ground vibration from blasting is measured as Peak Particle Velocity (PPV) in mm/s. Regulatory limits vary: typically 5–10 mm/s at nearest residence, 25–50 mm/s at mine infrastructure. Vibration is controlled by maximum instantaneous charge (MIC) and distance. The scaled distance formula (PPV = K × (D/Q^0.5)^-B) is the standard prediction tool.",
			},
		},
	})

	createSlide(database, s5.ID, "Misfire Management", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Misfire Management",
		"introText": "A misfire is any hole or portion of a hole that fails to detonate as designed. Misfires are the most dangerous post-blast condition — unexploded explosives buried in the muckpile can detonate during loading, hauling, or crushing. Misfire management is governed by regulation and must never be improvised.",
		"items": []map[string]interface{}{
			{
				"title":       "Types of Misfire",
				"description": "Complete misfire (entire hole fails to detonate), partial misfire (part of the column detonates, leaving unexploded explosive), and cut-off (detonation signal interrupted between holes, leaving downstream holes unfired). Electronic detonators significantly reduce cut-off risk because each detonator fires independently.",
			},
			{
				"title":       "Detection",
				"description": "Misfires are detected by: visual inspection of the muckpile (unbroken ground, intact stemming, visible detonator leads), comparison of post-blast survey to pre-blast design, and electronic detonator system logs (which report each unit's firing status). Post-blast walk-over by the shotfirer is mandatory.",
			},
			{
				"title":       "Treatment",
				"description": "A confirmed misfire requires isolation of the affected area, re-drilling and re-shooting the unfired hole (if safe to approach), or careful excavation under shotfirer supervision. The specific protocol is jurisdiction-dependent. In no case should a misfire area be loaded or hauled without shotfirer clearance.",
			},
			{
				"title":       "FMS and Load & Haul Integration",
				"description": "The misfire location must be communicated to the FMS and the excavator operator. The FMS can create a geofenced exclusion zone around the misfire location. The excavator must not dig within the exclusion zone until the shotfirer clears it. This links directly to the Fleet Management Systems and Load & Haul modules.",
			},
		},
	})

	createSlide(database, s5.ID, "Regulatory Requirements", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Regulatory Requirements",
		"introText": "Drill and blast operations are among the most heavily regulated activities in mining. The regulatory framework covers licensing, storage, transport, use, and disposal of explosives, as well as environmental controls on vibration, airblast, and flyrock.",
		"items": []map[string]interface{}{
			{
				"title":       "Shotfirer Licensing",
				"description": "Only a licensed shotfirer may initiate a blast. Licensing requirements vary by jurisdiction but typically require: formal training and examination, supervised practical experience (100+ blasts), first aid certification, and periodic renewal. The shotfirer is personally and legally responsible for blast safety.",
			},
			{
				"title":       "Explosives Storage and Handling",
				"description": "Explosives magazines must comply with AS 2187.1 (Australia) or equivalent jurisdiction standard. Requirements include: separation distances from other magazines and occupied buildings, security (locked, alarmed, access-controlled), temperature and humidity control, stock reconciliation (every unit accounted for), and transport in approved vehicles with dangerous goods placarding.",
			},
			{
				"title":       "Environmental Limits",
				"description": "Vibration (PPV) and airblast (dB(L)) limits are set by the mining licence or environmental approval. Exceedance triggers investigation, potential production curtailment, and regulatory action. Continuous blast monitoring (seismographs at compliance locations) is standard. Results must be recorded and available for regulatory audit.",
			},
		},
	})

	createSlide(database, s5.ID, "Misfire Procedures", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Misfire Procedures",
		"introText": "A misfire is any hole or charge that fails to detonate as planned. Misfires are the single most dangerous post-blast condition — the procedure is regulated, specific, and must never be improvised.",
		"items": []map[string]interface{}{
			{"title": "Definition", "description": "Any hole or portion of a charge that fails to detonate at the designed time. Includes complete misfires, partial column detonations, and cut-offs where downstream holes did not receive the firing signal."},
			{"title": "Immediate Action — Stay Out", "description": "The blast exclusion zone remains fully in force. Do NOT re-enter. No personnel, no equipment, no shortcuts. The unexploded charge can detonate with any disturbance."},
			{"title": "Notification Chain", "description": "Shot firer is first on scene (the licensed blaster who fired the shot). Shot firer reports to blast supervisor. Blast supervisor escalates to site manager. All three must know before any re-entry decision."},
			{"title": "Mandatory Wait Times", "description": "Electronic detonators: minimum 30 minutes before any re-entry (allows residual charge to stabilise). Shock tube (non-electric): minimum 1 hour. Some jurisdictions require longer — follow site procedure if stricter."},
			{"title": "Visual Inspection", "description": "From a safe distance only, ideally with binoculars or drone. Look for undetonated holes (intact stemming), surface expression of unfired charge, visible detonator leads, and any unusual muckpile geometry."},
			{"title": "Re-Entry", "description": "Only the trained shot firer, following the site-specific misfire protocol. No production personnel. Approach on foot with explosive detection equipment. The shot firer decides on re-shooting, manual extraction, or leave-in-place with exclusion geofence."},
			{"title": "Wet-Hole Misfires", "description": "Exceptionally high risk. Water desensitises emulsion and ANFO in unpredictable ways — the charge can be dormant for hours then detonate spontaneously. Do NOT attempt to extract. Specialist contractor treatment only."},
			{"title": "Regulatory Documentation", "description": "Every misfire is a reportable event. Documented in the shot firer's log, blast record, and site regulatory register. Root cause investigation is required before the next blast is fired."},
		},
	})
	createSlide(database, s5.ID, "Environmental & Infrastructure Blast Limits", "CONTENT", 5, map[string]interface{}{
		"type":      "content",
		"heading":   "Environmental & Infrastructure Blast Limits",
		"introText": "Every blast is bounded by regulatory limits on what the blast can do beyond the fence line. Flyrock, ground vibration, and airblast overpressure each have hard numerical limits — exceedance triggers investigation, production curtailment, and potential licence action.",
		"items": []map[string]interface{}{
			{"title": "Flyrock Exclusion Zone", "description": "Minimum 500 m from the blast perimeter in all directions for surface blasts. Site-specific risk assessments may extend this — large blasts, adverse geometry, or nearby roads often require 800 m or more. Exclusion is enforced by physical sentries, not signs."},
			{"title": "Ground Vibration — PPV", "description": "Peak Particle Velocity measured in mm/s at the nearest structure. Typical limit: <10 mm/s at occupied structures. Sensitive structures (heritage buildings, hospitals, precision infrastructure): <2 mm/s. Exceedance = investigation and report to regulator."},
			{"title": "Airblast Overpressure", "description": "Typical limit: <115 dB(L) at sensitive receptors (residences, schools). Airblast is the dominant complaint driver even when within limits — community perception is often ahead of regulatory thresholds."},
			{"title": "Monitoring", "description": "Seismograph (triaxial geophone + airblast microphone) at the nearest compliance location for every blast. Data logged, timestamped, and available for regulatory audit. Modern stations are solar-powered, cellular-connected, and can be deployed in 30 minutes."},
			{"title": "Jurisdictional Variation", "description": "Australia — each state differs. NSW: AS 2187.2 framework, typical 10 mm/s / 115 dB(L) at residences. QLD: similar numerical limits, different monitoring regime. WA: DMIRS guidelines, commonly tighter limits near population centres. International sites follow host-country standards. Know the specific licence limits before the first blast."},
			{"title": "Pre-Blast Notifications", "description": "Neighbouring residents, infrastructure owners (power, rail, pipeline operators), and aviation authorities (if near flight paths) must be notified in advance. Typical: 24-hour written notice plus same-day confirmation. Failure to notify = licence breach independent of whether any limit was exceeded."},
		},
	})
	createSlide(database, s5.ID, "Electronic Detonator Programming", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Electronic Detonator Programming",
		"introText": "Electronic detonators have replaced shock tube as the standard for high-value production blasts. Each detonator is individually addressable and millisecond-programmable — the blast design is literally uploaded to the pattern, hole by hole.",
		"items": []map[string]interface{}{
			{"title": "What They Replace", "description": "Shock tube (Nonel) uses pre-manufactured pyrotechnic delays in fixed increments (25, 42, 65 ms etc.) with ±5–10 ms scatter. Electronic detonators (e-dets) use a programmable digital timer with 1 ms resolution and <0.5 ms scatter."},
			{"title": "Unique Serial Number", "description": "Every e-det carries a unique serial number (and cryptographic ID on modern systems). The number is read into the blast box / logger when the detonator is deployed — no two devices can be confused, and the system records exactly which serial went into which hole."},
			{"title": "Programming Workflow", "description": "Step 1: tagger walks the pattern with a handheld logger and scans each detonator's ID at its hole location. Step 2: the logger auto-assigns the designed delay time for that hole from the blast design file. Step 3: the complete logged network is uploaded to the blast box (the firing controller)."},
			{"title": "Network Test", "description": "Before arming, the blast box runs a continuity check across every detonator in the network. The system reports any unit that did not respond — open circuits, damaged wires, failed electronics. The test must pass 100% before arming is permitted."},
			{"title": "Benefits Over Shock Tube", "description": "Millisecond precision timing (enables vibration-managed designs and fragmentation optimisation), individual unit verification (misfires detected pre-blast, not post-blast), and programmable firing sequences impossible with fixed delays (curved fronts, echelon, variable inter-hole timing)."},
			{"title": "Arming", "description": "Only after the exclusion zone has been physically confirmed clear and the shot firer has authorised firing. Arming energises the detonators; once armed, the system is committed. A separate fire command then triggers the sequence."},
			{"title": "Common Suppliers", "description": "Orica i-kon (most common in Australian open-cut), Dyno Nobel DigiShot (also widely used), Maxam MAXDET (strong presence in South America and Iberia). All systems are functionally similar but not interchangeable — loggers and blast boxes are product-specific."},
		},
	})

	// ── Section 6: Post-Blast Assessment ─────────────────────────────────
	s6 := createSection(database, module.ID, "Post-Blast Assessment", 5)

	createSlide(database, s6.ID, "Fragmentation Analysis", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Fragmentation Analysis",
		"introText": "Fragmentation analysis is how you close the feedback loop between blast design and blast outcome. Without measuring what you produced, you cannot know whether the design achieved its intent — and you cannot improve the next blast. The goal is not perfect fragmentation prediction, but systematic measurement that enables iterative improvement.",
		"items": []map[string]interface{}{
			{
				"title":       "Photo Fragmentation Analysis",
				"description": "Digital image analysis systems (WipFrag, Split-Desktop, PortaMetrics) photograph the muckpile surface and use edge-detection algorithms to estimate the size distribution. Modern systems use fixed cameras on excavator booms or drones to capture images automatically. Limitation: surface measurement only — the fines fraction below the surface is not captured, leading to systematic overestimation of mean fragment size.",
			},
			{
				"title":       "Drone Photogrammetry",
				"description": "Drones flying a grid pattern over the muckpile generate high-resolution orthophotos and 3D surface models. The 3D model provides muckpile volume estimation and surface fragmentation analysis across the entire blast, not just a single point. Drone surveys take 15–30 minutes and can be performed immediately post-blast.",
			},
			{
				"title":       "Physical Sieving",
				"description": "The gold standard for fragmentation measurement but expensive and slow. A representative sample is excavated, transported, and passed through a series of screens. Sieving provides a true size distribution including fines. Typically used for calibration of image analysis systems rather than routine measurement.",
			},
			{
				"title":       "Fragmentation Distribution Metrics",
				"description": "P80 (80% passing size) is the most common single-number descriptor. Typical targets: P80 = 300–600 mm for load and haul optimisation, P80 < 200 mm for direct crusher feed. The full distribution curve (% passing vs size) reveals whether the blast produced uniform fragmentation or a bimodal distribution (fines + oversized with a gap in the middle).",
			},
		},
	})

	createSlide(database, s6.ID, "Good vs Poor Fragmentation", "IMAGE", 1, map[string]interface{}{
		"type":        "image",
		"heading":     "Fragmentation Comparison — Well-Designed vs Under-Charged Blast",
		"description": "Side-by-side comparison of muckpile fragmentation. Left: a well-designed blast producing consistent, uniformly sized fragments (P80 ~400 mm) that an excavator can load efficiently. Right: an under-charged blast producing a mix of fine material and large boulders (>1 m), requiring secondary breakage and causing bucket stalling, reduced fill factors, and crusher bridging.",
		"imageUrl":    "https://assets.miningacademy.com/images/fragmentation-good-vs-poor.jpg",
	})

	createSlide(database, s6.ID, "Muckpile Profile and Dig Rate", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Muckpile Profile and Dig Rate",
		"introText": "Beyond fragment size, the physical shape and looseness of the muckpile determine how productively the excavator can work. A well-heaved, loose muckpile allows the bucket to penetrate easily and fill in 3–4 passes. A tight, poorly fragmented muckpile requires the excavator to crowd the face, increasing cycle times and wear.",
		"items": []map[string]interface{}{
			{
				"title":       "Muckpile Swell",
				"description": "Good blasts produce 30–40% swell (the muckpile volume is 30–40% larger than the in-situ rock volume). This indicates the rock has been adequately broken and loosened. Poor swell (<20%) suggests insufficient energy — the rock is cracked but not displaced.",
			},
			{
				"title":       "Muckpile Shape",
				"description": "The ideal muckpile for loading is a uniform, moderately heaved pile with a consistent height profile. Timing design controls this: V-patterns create a compact central pile, echelon patterns spread material laterally. The excavator's position relative to the muckpile must be considered in the blast design.",
			},
			{
				"title":       "Dig Rate Impact",
				"description": "Excavator dig rate (tonnes per operating hour) is directly correlated with fragmentation quality. Industry data shows: well-fragmented rock = 4,000–6,000 tph for a 400t class excavator; poorly fragmented rock = 2,500–3,500 tph for the same machine. The difference is 30–50% productivity.",
			},
			{
				"title":       "Toe Conditions",
				"description": "The bench floor after blasting should be clean and at grade. Poor sub-drill or insufficient energy at the toe leaves unbroken material (toe problems) that the excavator cannot dig. Toe problems require re-drilling and secondary blasting — a significant production delay. Consistent sub-drill (0.3 × burden) prevents this.",
			},
		},
	})

	createSlide(database, s6.ID, "The Blast Improvement Cycle", "DIAGRAM", 3, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Blast Improvement Cycle",
		"description": "Systematic blast improvement requires a closed feedback loop: design, execute, measure, compare, adjust. Without measurement, every blast is a guess. The sites that achieve best-in-class fragmentation performance run this cycle on every blast.",
		"diagramType": "cycle",
		"diagramData": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"label": "Design", "description": "Set pattern geometry, explosives, timing based on rock model and target fragmentation"},
				{"label": "Execute", "description": "Drill, load, fire per design. Record all as-built variations (hole depths, charge weights, water conditions)"},
				{"label": "Measure", "description": "Drone survey, photo fragmentation, dig rate monitoring, vibration/airblast records"},
				{"label": "Compare", "description": "Compare actual fragmentation, muckpile profile, and dig rate to design predictions and targets"},
				{"label": "Adjust", "description": "Modify pattern geometry, powder factor, timing, or explosive product for the next blast based on measured gaps"},
			},
		},
	})

	// ── Section 7: Mine to Mill — Integration ────────────────────────────
	s7 := createSection(database, module.ID, "Mine to Mill — Integration with the Value Chain", 6)

	createSlide(database, s7.ID, "The Mine to Mill Concept", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "The Mine to Mill Concept",
		"introText": "Mine to Mill is the principle that the entire value chain — from drill and blast through to final product — should be optimised as a single system, not as independent cost centres. Optimising drill and blast in isolation (minimising blast cost) sub-optimises the total system. The operational sites that achieve the lowest total cost per tonne deliberately over-invest in blast quality to reduce downstream costs.",
		"items": []map[string]interface{}{
			{
				"title":       "The Silo Problem",
				"description": "In a traditional mine, the drill and blast team minimises blast cost (powder factor, detonator spend), the load and haul team minimises fleet cost, and the processing team minimises mill energy. Each silo meets its own KPI while the total system underperforms. The drill and blast team has no incentive to spend more on explosives to help the crusher.",
			},
			{
				"title":       "Total Cost Optimisation",
				"description": "A 20% increase in powder factor might cost $0.10/t more in explosive but reduce crusher energy by $0.15/t and increase mill throughput by 5% (worth $0.30–0.50/t). The total system saves $0.35–0.55/t. This requires cross-functional KPIs and management systems that reward total cost reduction, not departmental cost minimisation.",
			},
			{
				"title":       "Fragmentation as the Control Variable",
				"description": "The link between drill and blast and downstream performance is fragmentation. P80 is the handover metric: the blast engineer designs for a target P80, the load and haul team reports actual dig rate vs P80, and the crusher reports throughput vs feed size. When P80 is measured and shared across departments, the mine to mill loop closes.",
			},
			{
				"title":       "Real-World Results",
				"description": "Published mine to mill case studies consistently show 10–20% total cost reduction when blast fragmentation is optimised for downstream performance. Notable examples: BHP Olympic Dam achieved 12% mill throughput increase through blast optimisation; Newcrest Cadia achieved 15% reduction in SAG mill energy consumption by increasing blast powder factor.",
			},
		},
	})

	createSlide(database, s7.ID, "How Blast Quality Cascades Through the Value Chain", "DIAGRAM", 1, map[string]interface{}{
		"type":        "diagram",
		"heading":     "How Blast Quality Cascades Through the Value Chain",
		"description": "Every downstream process is affected by what happens at the blast. This diagram traces the cascade of a fragmentation decision through loading, hauling, crushing, and milling. The numbers are representative of industry benchmarks.",
		"diagramType": "flow",
		"diagramData": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"label": "Blast: P80 = 400 mm", "description": "Well-fragmented. PF = 0.55 kg/m³. Blast cost: $0.45/t"},
				{"label": "Load: 4,800 tph", "description": "Good bucket fill, 4-pass loading. Dig rate benefit: +25%"},
				{"label": "Haul: 340t payload", "description": "Consistent bucket fills = consistent payloads. +8% vs poor blast"},
				{"label": "Crush: 3,200 tph", "description": "Feed size within crusher acceptance. No bridging. Throughput at design."},
				{"label": "Mill: 18.5 kWh/t", "description": "Finer feed = lower SAG mill energy. –12% vs coarse feed."},
				{"label": "Total: $4.80/t", "description": "vs $5.60/t for a low-PF blast that saved $0.15/t on explosives."},
			},
		},
	})

	createSlide(database, s7.ID, "Drill & Blast KPIs", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Drill & Blast KPIs",
		"introText": "These are the metrics that drive drill and blast performance and enable the mine to mill feedback loop. Each KPI connects blast quality to a measurable operational outcome.",
		"items": []map[string]interface{}{
			{
				"title":       "Powder Factor (kg/m³ or kg/t)",
				"description": "Primary design parameter. Target set by mine to mill optimisation, not minimum blast cost. Tracked per blast and trended over time. Deviation from design PF is a leading indicator of fragmentation variance.",
			},
			{
				"title":       "P80 Fragmentation (mm)",
				"description": "80% passing size measured by image analysis or physical sieving. The handover metric between drill and blast and load and haul. Tracked per blast and correlated with dig rate and crusher throughput.",
			},
			{
				"title":       "Drill Metres per Tonne (m/t)",
				"description": "Total drill metres divided by total tonnes broken. Reflects pattern density. Increasing drill metres per tonne (tighter patterns) costs more to drill but improves energy distribution and fragmentation uniformity.",
			},
			{
				"title":       "Cost per Tonne Broken ($/t)",
				"description": "Total drill and blast cost (drilling + explosives + accessories + labour) per tonne of rock broken. The direct cost metric. Must be viewed alongside downstream performance — not in isolation.",
			},
			{
				"title":       "Dig Rate (t/operating hour)",
				"description": "The downstream indicator of blast quality as experienced by the excavator. A sudden drop in dig rate at a face that was previously productive is a signal that the blast quality in that zone is below standard.",
			},
		},
	})

	createSlide(database, s7.ID, "Technology in Drill & Blast", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Technology in Drill & Blast",
		"introText": "Technology is transforming drill and blast from an experience-based craft to a data-driven engineering discipline. The common thread is measurement: better measurement of the rock, better measurement of the blast outcome, and faster feedback between the two.",
		"items": []map[string]interface{}{
			{
				"title":       "Measure While Drilling (MWD)",
				"description": "Drill rigs equipped with sensors record penetration rate, rotary pressure, pulldown force, and vibration for every metre of hole. These parameters correlate with rock hardness, fracture frequency, and water conditions. MWD data enables the blast design to be adjusted hole-by-hole based on actual rock conditions — not assumed from the geological model.",
			},
			{
				"title":       "Drone Surveying",
				"description": "Pre-blast: 3D bench surface model for accurate burden measurement and collar positioning. Post-blast: muckpile volume, surface fragmentation, and heave measurement. Turnaround time of 30–60 minutes from flight to data.",
			},
			{
				"title":       "Electronic Detonators",
				"description": "Programmable timing with 1 ms precision. Enables complex firing sequences, vibration management, and fragmentation optimisation that are impossible with pyrotechnic delays. The cost premium (3–5× per unit) is recovered in fragmentation improvement and vibration compliance.",
			},
			{
				"title":       "Fragmentation Cameras",
				"description": "Fixed cameras mounted on excavator booms or at crusher feeds capture images of every bucket or truck load. AI-based image analysis provides real-time P80 and size distribution, creating a continuous fragmentation record linked to blast location and design parameters.",
			},
			{
				"title":       "Blast Modelling Software",
				"description": "3D blast simulation (SHOTPlus, BlastMap, AEGIS) integrates pattern geometry, MWD rock characterisation, and timing design to predict fragmentation, vibration, and flyrock before the blast is fired. The simulation can be validated against post-blast measurement and calibrated over time.",
			},
		},
	})

	// ── Section 8: Autonomous and Digital Mines ──────────────────────────
	s8 := createSection(database, module.ID, "Drill & Blast in Autonomous and Digital Mines", 7)

	createSlide(database, s8.ID, "Autonomous Drilling", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Autonomous Drilling",
		"introText": "Autonomous drill systems are the most mature application of autonomy in drill and blast. Unlike autonomous haulage, which requires complex path planning and dynamic obstacle avoidance, drilling is a repetitive, pattern-based task well suited to automation. Major OEMs (Caterpillar, Epiroc, Sandvik) offer autonomous drill systems in production use.",
		"items": []map[string]interface{}{
			{
				"title":       "How Autonomous Drilling Works",
				"description": "The blast design is uploaded to the drill rig as a GPS-referenced pattern file. The drill navigates autonomously to each collar location, levels the mast, drills to the design depth, and moves to the next hole. A single operator supervises multiple rigs remotely. Collar positional accuracy: ±100 mm (significantly better than manual operator placement).",
			},
			{
				"title":       "Productivity Gains",
				"description": "Autonomous drilling eliminates operator-dependent variability: consistent collar positioning, consistent drill angle, no meal breaks or shift change delays. Published productivity gains: 15–30% increase in metres drilled per shift. The primary gain is utilisation (the drill keeps drilling during shift change, crib, and pre-start).",
			},
			{
				"title":       "MWD Integration",
				"description": "Autonomous drills capture continuous MWD data as a standard output. This data feeds directly into blast design software, enabling automated hole-by-hole charge design adjustment based on measured rock conditions. The link between autonomous drilling and adaptive blast design is the highest-value integration point in drill and blast technology.",
			},
		},
	})

	createSlide(database, s8.ID, "Digital Blast Management", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Digital Blast Management",
		"introText": "The digital mine integrates data from drilling, charging, firing, and post-blast assessment into a single digital record for every blast. This creates a dataset that enables machine learning, systematic improvement, and institutional knowledge that does not leave with individual blast engineers.",
		"items": []map[string]interface{}{
			{
				"title":       "Digital Blast Records",
				"description": "Every blast documented: as-designed pattern, as-drilled hole positions and depths (from autonomous drill GPS), MWD data, as-loaded charge weights and explosive products, detonator timing, post-blast fragmentation, vibration and airblast records, and excavator dig rates in the blasted zone. This is the blast's complete digital twin.",
			},
			{
				"title":       "Machine Learning Applications",
				"description": "With sufficient digital blast records (hundreds of blasts), ML models can predict fragmentation from design parameters and MWD data with accuracy that exceeds empirical models (Kuz-Ram). These models learn site-specific relationships between rock properties, explosive energy, and outcomes that empirical models cannot capture.",
			},
			{
				"title":       "Cross-Module Integration",
				"description": "The blast record links to the FMS (which tracks loading and hauling from the blasted zone), the crusher control system (which records throughput and energy consumption for the blasted material), and the grade control system (which classifies material). This integration is the digital implementation of mine to mill.",
			},
		},
	})

	createSlide(database, s8.ID, "Blast Monitoring Systems", "IMAGE", 2, map[string]interface{}{
		"type":        "image",
		"heading":     "Integrated Blast Monitoring — Seismographs and Airblast Sensors",
		"description": "A blast monitoring station positioned at the regulatory compliance point (nearest residence or infrastructure boundary). The station combines a triaxial geophone (measuring ground vibration in three axes as Peak Particle Velocity in mm/s) and an airblast microphone (measuring overpressure in dB(L)). Data is transmitted in real time to the blast engineer and automatically compared against licence limits. Modern stations are solar-powered, cellular-connected, and can be deployed in 30 minutes.",
		"imageUrl":    "https://assets.miningacademy.com/images/blast-monitoring-seismograph.jpg",
	})

	// ── Section 9: Assessment ────────────────────────────────────────────
	s9 := createSection(database, module.ID, "Assessment", 8)

	quizDesc := "Test your understanding of drill and blast operations — from drilling equipment and pattern design through explosives, blast execution safety, post-blast assessment, and the mine-to-mill concept."
	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Drill & Blast Operations Assessment",
		Description:      &quizDesc,
		PassMark:         80,
		MaxAttempts:      3,
		ShuffleQuestions: true,
		ShuffleOptions:   true,
		ShowResults:      true,
		SectionID:        &s9.ID,
		ModuleID:         &module.ID,
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
			question:    "What percentage of total mining cost does drill and blast typically represent, and what percentage of downstream cost does it influence?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Drill and blast is 15–20% of total mining cost but determines the efficiency of the 80–85% that follows — loading, hauling, crushing, and processing. This cost multiplier is why blast quality is the single highest-leverage investment in open-cut mining.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"5–10% of cost, influences 50% downstream", false},
				{"15–20% of cost, influences 80–85% downstream", true},
				{"40–60% of cost, influences 90% downstream", false},
				{"30–40% of cost, influences 60% downstream", false},
			},
		},
		{
			question:    "In what rock conditions should Down-the-Hole (DTH) drilling be preferred over rotary drilling?",
			qType:       models.QuestionMultipleChoice,
			explanation: "DTH drilling uses a pneumatic percussive hammer at the bottom of the hole, maintaining consistent energy transfer regardless of depth. It outperforms rotary in hard rock (UCS > 150 MPa) where rotary penetration drops sharply.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Soft sedimentary rock with UCS below 100 MPa", false},
				{"Hard, competent rock with UCS above 150 MPa", true},
				{"Any rock where large hole diameters (>270 mm) are required", false},
				{"Weathered material near the surface", false},
			},
		},
		{
			question:    "What is the recommended stemming height relative to burden?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Stemming height of 0.7–1.0 × burden provides adequate confinement to prevent gas venting through the collar. Too short = flyrock and airblast from the collar. Too tall = wasted hole length that could carry explosive.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"0.3 × burden", false},
				{"0.5 × burden", false},
				{"0.7–1.0 × burden", true},
				{"1.5 × burden", false},
			},
		},
		{
			question:    "A 200 mm diameter hole in a 5.0 m × 6.0 m pattern with 12 m bench height and ANFO (0.85 g/cc) produces a powder factor of approximately 0.45 kg/m³. The blast engineer wants to increase the powder factor to 0.55 kg/m³ without changing hole diameter. What is the most effective approach?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Reducing burden and spacing (tighter pattern) increases the number of holes per area, increasing explosive mass per cubic metre of rock. While switching to denser emulsion also increases PF, tighter patterns both increase PF and improve energy distribution uniformity.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Increase bench height to 15 m", false},
				{"Reduce burden and spacing (tighter pattern)", true},
				{"Reduce stemming height to 1.5 m", false},
				{"Switch to emulsion at 1.2 g/cc and keep the same pattern", false},
			},
		},
		{
			question:    "What is the primary advantage of electronic detonators over non-electric (Nonel) detonators?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Electronic detonators provide programmable delay with 1 ms precision and <0.5 ms scatter. This enables complex timing designs, precise vibration control, and fragmentation optimisation that are impossible with the fixed delay increments and ±5–10 ms scatter of Nonel systems.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Lower cost per unit", false},
				{"Faster hookup time", false},
				{"Programmable timing with 1 ms precision and near-zero scatter", true},
				{"Ability to detonate in wet conditions", false},
			},
		},
		{
			question:    "During blast clearance, what must be confirmed before the shotfirer can initiate the firing sequence?",
			qType:       models.QuestionMultipleChoice,
			explanation: "The shotfirer must personally confirm: sentries at all access points, personnel and equipment count reconciled to zero inside the exclusion zone, and the standard siren sequence completed. Only the licensed shotfirer can initiate the blast.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"All trucks have left the pit", false},
				{"All sentries are in position, all personnel accounted for, and exclusion zone is secure", true},
				{"The FMS has automatically cleared all vehicles from the blast zone", false},
				{"The supervisor has given verbal approval", false},
			},
		},
		{
			question:    "What are the three mechanisms that cause flyrock?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Flyrock is propelled by blast gases through three mechanisms: cratering (gas venting through insufficient stemming at the collar), face burst (gas breaking through the face due to excessive burden), and rifling (gas following joints or weak zones, ejecting material at unexpected trajectories).",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Vibration, airblast, and ground heave", false},
				{"Cratering (collar venting), face burst (inadequate free face), and rifling (gas following weak zones)", true},
				{"Overloading, poor timing, and wet holes", false},
				{"Excessive sub-drill, high powder factor, and large hole diameter", false},
			},
		},
		{
			question:    "A misfire has been identified in a muckpile that is about to be loaded. What is the correct action?",
			qType:       models.QuestionMultipleChoice,
			explanation: "A misfire means unexploded explosives are buried in the muckpile. The area must be isolated, the shotfirer notified, and no loading or hauling may occur until the shotfirer has treated the misfire and cleared the zone. There is no safe distance for digging near a misfire without shotfirer supervision.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Continue loading but avoid the misfire area by 10 m", false},
				{"The excavator operator may dig around it if they proceed carefully", false},
				{"Isolate the area, notify the shotfirer, and do not load or haul until the shotfirer has cleared the zone", true},
				{"Mark it on the FMS and continue loading from the other end", false},
			},
		},
		{
			question:    "What is the primary limitation of photo fragmentation analysis (image analysis)?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Image analysis measures the surface of the muckpile. Finer material segregates to the bottom and interior, so the surface is biased toward coarser fragments. This systematically overestimates mean fragment size and underestimates the fines fraction.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"It cannot measure fragments larger than 500 mm", false},
				{"It is a surface measurement only — it overestimates mean size by missing fines below the surface", true},
				{"It requires the muckpile to be completely dry", false},
				{"It can only be performed from the ground, not from drones", false},
			},
		},
		{
			question:    "The Mine to Mill concept argues that increasing powder factor can reduce total mining cost. What is the mechanism?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Increasing PF costs more in explosive but produces finer, more uniform fragmentation. This improves excavator dig rate (+25–30%), truck payload consistency (+5–10%), crusher throughput, and reduces mill energy consumption (–10–15%). The downstream savings typically exceed the explosive cost increase by 3–5×.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Higher powder factor reduces the amount of drilling required", false},
				{"Finer fragmentation from higher PF improves load, haul, crush, and mill performance — saving more downstream than the extra explosive costs", true},
				{"Higher powder factor creates more swell, which makes trucks look fuller", false},
				{"The blast team receives a bonus for higher powder factor, incentivising better design", false},
			},
		},
		{
			question:    "What data does Measure While Drilling (MWD) provide, and how is it used in blast design?",
			qType:       models.QuestionMultipleChoice,
			explanation: "MWD sensors record penetration rate, rotary pressure, pulldown force, and vibration continuously during drilling. These parameters correlate with rock hardness and fracture frequency, enabling the blast design to be adjusted hole-by-hole based on actual rock conditions rather than assumed geology.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Hole diameter and depth measurements used to calculate explosive volume", false},
				{"Penetration rate, rotary pressure, and pulldown force that characterise rock hardness for hole-by-hole charge adjustment", true},
				{"GPS coordinates of each hole used to generate the pattern map", false},
				{"Water table depth measurements used to select between ANFO and emulsion", false},
			},
		},
		{
			question:    "A blast has been fired and the post-blast survey shows P80 of 650 mm against a target of 400 mm. Dig rates in the muckpile are 30% below the previous blast. What is the most likely cause and the correct response?",
			qType:       models.QuestionMultipleChoice,
			explanation: "P80 of 650 mm vs target 400 mm indicates significantly coarser fragmentation than designed. This directly causes the 30% dig rate reduction. The response is to investigate the root cause (burden accuracy, powder factor, timing, rock conditions) using the blast feedback loop and adjust the next blast design.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"The excavator is underperforming — inspect the machine", false},
				{"Fragmentation is coarser than design — review burden/spacing accuracy, powder factor, and timing to identify the cause, and adjust the next blast design", true},
				{"The trucks are under-loaded — adjust the FMS payload targets", false},
				{"The post-blast survey is inaccurate — re-fly the drone", false},
			},
		},
	}

	for i, q := range questions {
		question := models.QuizQuestion{
			ID:          uuid.New().String(),
			QuizID:      quiz.ID,
			Question:    q.question,
			Type:        q.qType,
			Explanation: &q.explanation,
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

	createSlide(database, s9.ID, "Drill & Blast Operations Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Drill & Blast Operations Assessment",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{
				"question":     "What percentage of total mining cost does drill and blast represent, and what percentage of downstream cost does it influence?",
				"options":      []string{"5–10% of cost, influences 50% downstream", "15–20% of cost, influences 80–85% downstream", "40–60% of cost, influences 90% downstream", "30–40% of cost, influences 60% downstream"},
				"correctIndex": 1,
				"explanation":  "Drill and blast is 15–20% of total mining cost but determines the efficiency of the 80–85% that follows.",
			},
			{
				"question":     "What are the three mechanisms that cause flyrock?",
				"options":      []string{"Vibration, airblast, and ground heave", "Cratering (collar venting), face burst (inadequate free face), and rifling (gas following weak zones)", "Overloading, poor timing, and wet holes", "Excessive sub-drill, high powder factor, and large hole diameter"},
				"correctIndex": 1,
				"explanation":  "Flyrock is propelled by blast gases through cratering (collar), face burst (face), and rifling (weak zones).",
			},
			{
				"question":     "The Mine to Mill concept argues that increasing powder factor can reduce total mining cost. What is the mechanism?",
				"options":      []string{"Higher powder factor reduces the amount of drilling required", "Finer fragmentation improves load, haul, crush, and mill performance — saving more downstream than the extra explosive costs", "Higher powder factor creates more swell", "The blast team receives a bonus for higher powder factor"},
				"correctIndex": 1,
				"explanation":  "Finer fragmentation from higher PF saves 3–5× the additional explosive cost across downstream operations.",
			},
		},
	})

	// ── Section 10: Completion ──────────────────────────────────────────
	s10 := createSection(database, module.ID, "Completion", 9)

	createSlide(database, s10.ID, "Module Complete", "COMPLETION", 0, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "You have completed the Drill & Blast Operations module.",
		"summary": "You can now explain how drill and blast quality cascades through the entire mining value chain, describe drill pattern design parameters and their consequences, select the appropriate explosive and detonator for given conditions, follow blast clearance procedures and identify flyrock and misfire hazards, assess post-blast fragmentation and use the blast improvement cycle, and apply the mine-to-mill concept to optimise total cost rather than blast cost alone.",
	})

	log.Println("Drill & Blast Operations module seeded successfully.")
}
