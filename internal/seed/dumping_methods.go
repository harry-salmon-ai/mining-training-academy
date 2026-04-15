package seed

import (
	"log"

	"go-backend-react-frontend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDumpingMethodsModule(database *gorm.DB, force bool) {
	var moduleCount int64
	database.Model(&models.Module{}).Where("slug = ?", "dumping-methods").Count(&moduleCount)
	if moduleCount > 0 {
		if !force {
			log.Println("Dump Operations module already exists, skipping")
			return
		}
		forceDeleteModule(database, "dumping-methods")
	}

	log.Println("Seeding Dump Operations module...")

	var category models.Category
	database.Where("slug = ?", "mining-operations").First(&category)
	var admin models.User
	database.Where("role = ?", models.RoleSuperAdmin).First(&admin)

	if category.ID == "" || admin.ID == "" {
		log.Println("Cannot seed Dump Operations: category or admin not found")
		return
	}

	desc := "A practical guide to dump operations in open-cut mining — covering standard dumping procedures, highwall and OTE mechanics, capping dump technique, paddock and crusher dumps, road maintenance dumping, and how autonomous haulage systems change the dump environment."
	dur := 90
	module := models.Module{
		ID:          uuid.New().String(),
		Title:       "Dump Operations",
		Slug:        "dumping-methods",
		Description: &desc,
		Status:      models.StatusPublished,
		Level:       models.LevelFoundation,
		Duration:    &dur,
		CategoryID:  category.ID,
		AuthorID:    admin.ID,
	}
	database.Create(&module)

	tags := []string{"dump operations", "highwall dump", "OTE", "over-the-edge", "capping dump", "paddock dump", "crusher dump", "autonomous dumping", "AHS", "haul truck", "dump berm"}
	for _, t := range tags {
		database.Create(&models.ModuleTag{ID: uuid.New().String(), ModuleID: module.ID, Tag: t})
	}

	// ── Section 1: Introduction ─────────────────────────────────────────
	s1 := createSection(database, module.ID, "Introduction", 0)

	createSlide(database, s1.ID, "Dump Operations", "TITLE", 0, map[string]interface{}{
		"type":     "title",
		"heading":  "Dump Operations",
		"subtitle": "Technique, safety, and the advancing highwall",
		"bullets": []string{
			"Standard dump procedures every haul truck operator must know",
			"Highwall dumps — OTE mechanics and the capping cycle",
			"Paddock, crusher, and road maintenance dump methods",
			"How autonomous haulage changes the dump environment",
		},
	})

	createSlide(database, s1.ID, "Dumps in the Mining Value Chain", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Dumps in the Mining Value Chain",
		"introText": "Every haul truck cycle ends at a dump. Where material lands — and how it is placed — directly affects safety, geotechnical stability, production continuity, and the ability to advance the mine plan. Poor dump practice is one of the highest-consequence risks in open-cut mining.",
		"items": []map[string]interface{}{
			{
				"title":       "Safety Consequences",
				"description": "Dump failures — material slip, truck rollover, berm collapse — are typically high-energy events. The consequences of dumping onto unstable ground or without a berm range from equipment loss to fatalities.",
			},
			{
				"title":       "Production Continuity",
				"description": "Dump delays ripple back through the fleet. If the dump point is blocked, crushed, or under inspection, every truck behind it queues. Dump efficiency is a direct driver of cycle time and tonnes per hour.",
			},
			{
				"title":       "Mine Plan Execution",
				"description": "Waste dumps must be built according to the geotechnical plan — correct lift heights, drainage provisions, and crest positions. Incorrect dump placement can sterilise future pit access or create unstable faces that require expensive remediation.",
			},
		},
	})

	createSlide(database, s1.ID, "Dump Type Classification", "DIAGRAM", 2, map[string]interface{}{
		"type":        "diagram",
		"heading":     "Dump Type Classification",
		"description": "Dump type is determined by the material being moved, the receiving location, and site geometry. Understanding which dump type applies to each load is the first decision every operator makes at the end of a haul.",
		"diagramType": "tree",
		"diagramData": map[string]interface{}{
			"nodes": []map[string]interface{}{
				{"id": "root", "label": "Material Assignment", "description": "Every haul truck load is assigned to a dump based on material type and location"},
				{"id": "waste", "label": "Waste Material", "parent": "root", "description": "Overburden and non-ore material to be permanently placed"},
				{"id": "ore", "label": "ROM Ore", "parent": "root", "description": "Run-of-mine ore for processing"},
				{"id": "maint", "label": "Maintenance Material", "parent": "root", "description": "Gravel, crushed rock for road upkeep"},
				{"id": "hw", "label": "Highwall Dump", "parent": "waste", "description": "Active highwall face — OTE and capping technique required"},
				{"id": "pad", "label": "Paddock Dump", "parent": "waste", "description": "Open ground dump — layered lift construction"},
				{"id": "crusher", "label": "Crusher / ROM Pad", "parent": "ore", "description": "Direct delivery to primary crusher pocket or ROM stockpile"},
				{"id": "road", "label": "Road Maintenance Dump", "parent": "maint", "description": "Material placed on haul road for grading and compaction"},
			},
		},
	})

	// ── Section 2: Standard Dumping Procedure ──────────────────────────
	s2 := createSection(database, module.ID, "Standard Dumping Procedure", 1)

	createSlide(database, s2.ID, "Dump Point Layout from Above", "IMAGE", 0, map[string]interface{}{
		"type":        "image",
		"heading":     "Dump Point Layout from Above",
		"description": "Before executing the dump procedure, understand the spatial relationships at the dump face: the berm runs parallel to the edge, the truck reverses perpendicular to the face, and a minimum 5 m clear zone separates the berm from the crest. Every step in the universal dump procedure maps to one of these spatial zones.",
		"imageUrl":    "https://assets.miningacademy.com/images/dump-point-aerial-layout.jpg",
	})
	createSlide(database, s2.ID, "The Universal Dump Procedure", "DIAGRAM", 1, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Universal Dump Procedure",
		"description": "Every dump — regardless of type — follows this sequence. Skipping steps is the most common cause of dump incidents. The procedure applies to both manned and autonomous operations, though AHS executes it via system logic rather than operator judgment.",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"title": "Approach at controlled speed", "description": "Reduce to walk speed within 50 m of the dump face. Assess ground conditions, visibility, and any changes since the last cycle."},
				{"title": "Confirm dump assignment", "description": "Verify you are at the correct dump point via FMS dispatch, signage, or spotter instruction. Wrong-dump events are more common than expected."},
				{"title": "Inspect the berm", "description": "A continuous berm of at least half the rear tyre diameter must be present across the full tip width before reversing. No berm — no dump."},
				{"title": "Establish spotter communication", "description": "Where a spotter is required (highwall, crusher pocket, restricted visibility), make two-way contact before beginning the reverse. Confirm signals or radio channel."},
				{"title": "Reverse to the dump position", "description": "Reverse slowly and steadily, wheels aligned with the dump direction. Stop when rear tyres are 1–2 m from the berm or at the distance directed by the spotter."},
				{"title": "Raise the body and tip", "description": "Apply park brake, confirm flat ground (or within safe grade), raise tray. Observe for unexpected slumping, cracking, or movement during tip."},
				{"title": "Lower body and clear the dump", "description": "Lower tray fully before driving away — never travel with the body raised. Move forward at least one truck length before turning."},
				{"title": "Report any anomalies", "description": "Soft ground, berm failure, slumping, or unusual cracking must be reported to the supervisor before the next truck tips. Do not leave an unsafe dump unattended."},
			},
		},
	})

	createSlide(database, s2.ID, "Berm Requirements", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Berm Requirements",
		"introText": "The dump berm is the primary physical safeguard preventing a truck from going over the edge. Berm specifications are not suggestions — they are engineered controls backed by geotechnical design. The half-tyre-height rule is the minimum standard across Australian open-cut operations.",
		"items": []map[string]interface{}{
			{
				"title":       "The Half-Tyre-Height Rule",
				"description": "The berm height must be at least half the inflated diameter of the largest rear tyre operating at the dump. For a typical large haul truck with 4.0 m diameter tyres, this means a minimum 2.0 m berm height. Smaller berms allow a truck to roll over them under the weight of a loaded body.",
				"imageUrl":    "https://assets.miningacademy.com/images/berm-half-tyre-height-rule.jpg",
			},
			{
				"title":       "Windrow vs Constructed Berm",
				"description": "A windrow berm is formed by pushing previously tipped material into a continuous mound — fast to build but can degrade quickly with truck traffic and rain. A constructed berm uses compacted material of a specified grade. Highwall dump berms must be constructed and regularly inspected for gaps, erosion, and compaction.",
			},
			{
				"title":       "Berm Continuity",
				"description": "A berm with gaps is not a berm. The full width of the tip area must be covered without interruption. After every tip, the berm should be assessed for damage, and after large-body material tips (highwall OTE), the berm may need re-establishing before the next truck can use that section.",
			},
		},
		"keyFact": "No berm = no dump. This is non-negotiable. An operator who tips without a berm has bypassed the primary engineered safeguard at that dump point.",
	})

	createSlide(database, s2.ID, "Spotter Protocols", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Spotter Protocols",
		"introText": "Spotters are required whenever the driver cannot directly assess the dump face — restricted visibility, highwall edges, crusher pockets, and multi-truck dump points all trigger this requirement. A spotter is a last line of defence, not a replacement for pre-tip inspection.",
		"items": []map[string]interface{}{
			{
				"title":       "When Spotters Are Mandatory",
				"description": "Highwall dumps where the edge is not visible from the cab, crusher pocket dumps, any situation with less than 50 m clear visibility to the berm, night operations in restricted areas, and first use of a new dump point. Site procedures define the full list — when in doubt, use a spotter.",
			},
			{
				"title":       "Establishing Contact",
				"description": "Before any reverse movement begins, the driver must establish two-way communication with the spotter — either radio contact or confirmed visual signal exchange. The spotter acknowledges readiness and the driver confirms signal conventions (for hand signals: stop = closed fist, advance = beckoning, come back = rotating arm).",
			},
			{
				"title":       "The Stop Signal is Absolute",
				"description": "If a spotter gives a stop signal — by radio or by hand — the driver stops immediately, without question or delay. The driver does not move until the spotter re-establishes contact and gives a proceed instruction. Loss of spotter contact during a reverse is a stop condition.",
			},
		},
	})

	createSlide(database, s2.ID, "Stop Conditions", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Stop Conditions — When Not to Dump",
		"introText": "There are conditions under which no dump should proceed, regardless of dispatch instructions or production pressure. Recognising and acting on stop conditions is a critical operator competency. The cost of a refused dump is measured in minutes; the cost of ignoring a stop condition can be irreversible.",
		"items": []map[string]interface{}{
			{
				"title":       "No Berm or Inadequate Berm",
				"description": "If the berm is absent, too low, discontinuous, or has been disrupted by previous tips, the dump must not proceed. The operator should park clear of the dump area and report to the supervisor immediately.",
			},
			{
				"title":       "Ground Movement or Slumping",
				"description": "Cracks running parallel to the dump edge, slumping of previously tipped material, or soft ground underfoot are signs of instability. These must be reported and the dump point isolated. Never tip onto an actively moving or cracking face.",
			},
			{
				"title":       "Loss of Spotter Contact",
				"description": "If radio contact with the spotter is lost during approach or reversing, stop immediately and do not resume until contact is re-established. A silent radio is a stop condition — never assume 'no news is good news' when reversing blind.",
			},
			{
				"title":       "Poor Visibility",
				"description": "Dust, dark, or obscured lines of sight that prevent the driver from seeing the berm or receiving spotter signals are stop conditions. Wait for conditions to improve or deploy a spotter. Visibility does not improve with hesitation — it improves with patience.",
			},
		},
	})

	// ── Section 3: Highwall Dumps ───────────────────────────────────────
	s3 := createSection(database, module.ID, "Highwall Dumps", 2)

	createSlide(database, s3.ID, "Highwall Dump — Three Dimensions", "IMAGE", 0, map[string]interface{}{
		"type":        "image",
		"heading":     "Highwall Dump — Three Dimensions",
		"description": "A highwall dump is a three-dimensional structure built against the face of the pit. Material is tipped from the crest (top), builds a slope (face), and accumulates at the base (toe). The vertical scale — commonly 20–50 m of exposed face — is only apparent when you see it from inside the pit.",
		"imageUrl":    "https://assets.miningacademy.com/images/highwall-dump-face-from-pit.jpg",
	})
	createSlide(database, s3.ID, "What is a Highwall Dump?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What is a Highwall Dump?",
		"introText": "A highwall dump is an end-of-bench waste dump placed at the crest of a pit highwall. Material is tipped directly over the edge — the dump face is the highwall itself. Highwall dumps are space-efficient (waste goes directly into the mined void) but carry the highest geotechnical risk of any dump type.",
		"items": []map[string]interface{}{
			{
				"title":       "Why Use Highwall Dumps?",
				"description": "When waste dumps outside the pit are distant or unavailable, placing waste over the highwall minimises haul distance and total cycle time. On some operations, highwall dumping is the only economically viable method for managing large volumes of overburden quickly.",
			},
			{
				"title":       "The Highwall Geometry",
				"description": "The dump has three key elements: the crest (the edge from which trucks tip), the face (the highwall wall below the crest), and the toe (the material that accumulates at the base of the face). As dumps progress, the crest advances into previously unmined ground and the face height increases.",
			},
			{
				"title":       "Geotechnical Risk",
				"description": "Highwall dumps operate on an active edge over a significant drop. Crest failure can carry equipment and operators over the wall. Geotechnical engineers set safe crest positions, maximum lift heights, and inspection intervals — operators must follow these parameters exactly.",
			},
		},
		"keyFact": "Highwall dumps have the highest inherent risk of any dump type. The engineered controls — berm specification, crest setback, inspection frequency — are not conservative buffers. They are the minimum safe standard.",
	})

	createSlide(database, s3.ID, "The Dump Face Geometry", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "The Dump Face Geometry",
		"introText": "Understanding the physical structure of the highwall dump is essential before operating at one. The geometry changes with every load — which is why survey and boundary management are inseparable from highwall dump operations.",
		"items": []map[string]interface{}{
			{
				"title":       "Crest",
				"description": "The upper edge of the dump — where trucks position their rear tyres. The crest is defined by the berm. As dump cycles advance the face, the crest moves outward (away from the pit wall behind it). The crest position is surveyed and updated in the FMS after each significant advance.",
			},
			{
				"title":       "Face",
				"description": "The slope below the crest, running from the dump surface down to the toe. Face angle is a geotechnical design parameter — the steepness depends on material type, moisture, and the mine's geotechnical design. As OTE dumps accumulate, the slope builds and the toe gradually flattens as the material approaches its natural repose angle.",
			},
			{
				"title":       "Toe",
				"description": "The material that accumulates at the base of the face as OTE dumps build the slope. As the toe builds and the face flattens, the dump is approaching the point where a capping cycle is required to advance the face further. The toe condition is a key indicator of when to transition from OTE to capping mode.",
			},
		},
	})

	createSlide(database, s3.ID, "Over-the-Edge (OTE) Mechanics", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Over-the-Edge (OTE) Mechanics",
		"introText": "Over-the-edge dumping is the process by which waste material is placed directly over the highwall crest. The truck reverses to the berm, raises the body, and the material flows over the edge under gravity, rolling down the face to build the toe. OTE is the primary waste placement method at a highwall dump.",
		"items": []map[string]interface{}{
			{
				"title":       "The OTE Position",
				"description": "The truck's rear tyres stop against the berm — typically 1–2 m from the crest edge. The body is raised and material slides over the edge. The crest is not a fixed drop point; material runs down the face at the natural repose angle, and subsequent tips add to the same slope.",
			},
			{
				"title":       "Slope Buildup and Toe Flattening",
				"description": "Each OTE dump adds material to the slope. The slope builds progressively and the toe gradually flattens as material accumulates. After 3–4 OTE dumps per dump point, the face has built sufficiently for the capping cycle to begin — the flattening toe is the signal that the OTE phase is complete.",
			},
			{
				"title":       "OTE is Done Per Dump Point",
				"description": "A highwall dump face typically has multiple dump points along its length. Each point runs its OTE count independently. Once all points in the face have reached their OTE threshold, the entire face transitions to the capping phase together before the berm reconstruction can occur.",
			},
		},
	})

	createSlide(database, s3.ID, "OTE Dump in Action", "IMAGE", 4, map[string]interface{}{
		"type":        "image",
		"heading":     "OTE Dump in Action",
		"description": "An over-the-edge dump: the truck body is raised, material flows over the crest, and cascades down the highwall face under gravity. The slope building below the truck will become the foundation of the next OTE dump — and eventually the material the dozer pushes outward during the capping phase.",
		"imageUrl":    "https://assets.miningacademy.com/images/ote-dump-truck-tipping-highwall.jpg",
	})
	createSlide(database, s3.ID, "The OTE to Cap Cycle", "DIAGRAM", 5, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The OTE → Cap Cycle",
		"description": "Highwall dump advance follows a repeating cycle. Each dump point runs an OTE count; when all points reach threshold a capping dump is performed per point; the dump then closes while a dozer reconstructs the berm; then the cycle resets at the new advanced position.",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"title": "Dump point configured — OTE script active", "description": "Each highwall dump point is set up with an OTE dump script and a capping dump script. The FMS begins dispatching trucks under the OTE script and tracking the dump count per point."},
				{"title": "OTE dumps 1–N (typically 3–4 per point)", "description": "Trucks reverse to the berm and tip over the edge. Each dump builds the slope further; the toe gradually flattens as material accumulates toward its natural angle. The FMS increments the count with each tip."},
				{"title": "OTE count threshold reached — switch to cap script", "description": "When a dump point reaches its configured OTE count, the FMS switches that point to the capping dump script. The next truck dispatched to that point executes a capping dump rather than an OTE dump."},
				{"title": "Capping dump — material placed at crest", "description": "The capping truck tips its load at the dump point, depositing material at the crest. This provides the loose material the dozer will use to push the berm outward. Once all points in the face have capped, the dump closes."},
				{"title": "Dump closes — dozer reconstructs the berm", "description": "The FMS closes the highwall dump: no trucks are dispatched. A dozer operator works across all dump points, pushing the capping material outward to advance and reconstruct the berm along the full face."},
				{"title": "Dump reopens — boundary updated, count reset", "description": "Once the dozer has pushed out all points and the new berm is established, the dump boundary is updated in the FMS. The OTE count resets and trucks begin the next OTE cycle from the new advanced berm position."},
			},
		},
	})

	createSlide(database, s3.ID, "Highwall Dump Boundary Management", "CONTENT", 6, map[string]interface{}{
		"type":      "content",
		"heading":   "Highwall Dump Boundary Management",
		"introText": "The highwall dump berm advances with every OTE-to-cap cycle. The boundary in the FMS must reflect the physical berm position at all times — the dump is closed while the dozer works, and only reopens once the new boundary is set. An out-of-date boundary means trucks are sent to the wrong position.",
		"items": []map[string]interface{}{
			{
				"title":       "Boundary Advances After the Dozer, Not the Trucks",
				"description": "The capping dump places loose material at the crest. The dozer then pushes that material outward to physically form the new berm. Only after the dozer has completed the berm reconstruction is the new boundary position confirmed and loaded into the FMS. The trucks did not advance the berm — the dozer did.",
			},
			{
				"title":       "Dump Closure During Dozer Work",
				"description": "While the dozer is reconstructing the berm, no trucks access the dump. The FMS holds the dump in a closed state. This is not downtime — it is a required phase of the advance cycle. The closure window is predictable and can be planned into fleet dispatch to minimise queue buildup.",
			},
			{
				"title":       "Reopening the Dump",
				"description": "Once the dozer has pushed out all points along the face and the new berm is established, the updated boundary is loaded into the FMS and the dump is reopened. The OTE count resets to zero on each point. The next OTE cycle begins from the new, advanced berm position.",
			},
		},
	})

	// ── Section 4: Capping Dumps ────────────────────────────────────────
	s4 := createSection(database, module.ID, "Capping Dumps", 3)

	createSlide(database, s4.ID, "What is a Capping Dump?", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "What is a Capping Dump?",
		"introText": "A capping dump is a truck tip performed at a highwall dump point after its OTE count threshold has been reached. The truck tips its load at the dump point crest — not over the edge — providing loose material that a dozer will push outward to advance and reconstruct the berm. The capping dump is done once per dump point.",
		"items": []map[string]interface{}{
			{
				"title":       "The Purpose: Providing Material for the Dozer",
				"description": "The capping dump does not advance the berm itself — the dozer does. The truck's role is to deliver a load of material to the crest. The dozer then uses that material to push the berm outward to the new face position. Without the capping dump, the dozer has no material to work with.",
			},
			{
				"title":       "One Cap Per Dump Point",
				"description": "Each dump point along the highwall face receives one capping dump after its OTE threshold is met. Once all points have been capped, the dump closes and the dozer works across the full face — pushing out each point in sequence — to reconstruct the berm at the new advanced position.",
			},
			{
				"title":       "Capping Out a Point",
				"description": "\"Capping out a point\" means running the full OTE-to-cap-to-dozer cycle repeatedly until a dump point reaches its designed maximum extent. Each cycle advances the face one increment. When the dump boundary reaches the mine plan limit at a given point, that point is capped out and retired.",
			},
		},
	})

	createSlide(database, s4.ID, "The Capping Dump Sequence", "DIAGRAM", 1, map[string]interface{}{
		"type":        "diagram",
		"heading":     "The Capping Dump Sequence",
		"description": "Once the OTE count is complete across all dump points, the capping phase begins. The FMS switches each point to its cap script, collects one cap dump per point, then closes the dump for dozer work.",
		"diagramData": map[string]interface{}{
			"steps": []map[string]interface{}{
				{"title": "OTE threshold reached on all points", "description": "The FMS has tracked the OTE dump count at each point. Once every point in the face reaches its configured threshold, the transition to capping begins."},
				{"title": "FMS switches each point to cap script", "description": "The dump point configuration changes from OTE script to capping dump script. The reverse path and tip position may differ slightly from the OTE position — the cap is tipped at the crest, not over the edge."},
				{"title": "Cap truck tips at the crest", "description": "The truck reverses to the dump point and tips its load. The material lands at the crest — it is not tipped over the edge. This loose pile of material is what the dozer will push outward."},
				{"title": "All points capped — dump closes", "description": "Once every point in the face has received its capping dump, the FMS closes the entire dump. No trucks are dispatched. The dozer operator is given the task to work the face."},
				{"title": "Dozer pushes the berm outward", "description": "The dozer works across all dump points, pushing the capping material outward from each point to advance and reconstruct the berm along the full face. This is the step that physically advances the dump."},
				{"title": "New boundary set — dump reopens", "description": "With the berm rebuilt at the new position, the dump boundary is updated in the FMS. Each point's OTE count resets. The dump reopens under OTE script and the cycle begins again."},
			},
		},
	})

	createSlide(database, s4.ID, "Dozer Pushing the Berm Outward", "IMAGE", 2, map[string]interface{}{
		"type":        "image",
		"heading":     "Dozer Pushing the Berm Outward",
		"description": "After every capping dump cycle, a dozer works across the full highwall face — pushing the capping material outward from each dump point to advance and reconstruct the berm. This is the step that physically moves the dump face forward. The dozer completes all points before the dump reopens.",
		"imageUrl":    "https://assets.miningacademy.com/images/dozer-pushing-berm-highwall.jpg",
	})
	createSlide(database, s4.ID, "The Dozer's Role", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "The Dozer's Role",
		"introText": "The dozer is the agent that actually advances the highwall dump face. Trucks deliver material; the dozer shapes it into a safe, extended berm. Understanding the dozer's role clarifies why the dump must close — and why the closure window is a planned, unavoidable part of the advance cycle.",
		"items": []map[string]interface{}{
			{
				"title":       "Pushing the Berm Outward",
				"description": "The dozer takes the material left by each capping dump and pushes it forward to form the new berm line. This extends the safe tip zone further out along the face. Without the dozer push, the capping material would just sit at the old crest and provide no useful berm advance.",
			},
			{
				"title":       "Working Across All Points",
				"description": "The dozer works point-by-point along the full face. Each point gets its own berm push. The sequence and timing of the dozer work across points is an operational planning decision — typically the dozer works from one end of the face to the other in a single pass.",
			},
			{
				"title":       "The Dump is Closed While the Dozer Works",
				"description": "No trucks access the dump during dozer work. This is a hard requirement — trucks and a working dozer on the same active face is an interaction hazard. The dump closure is managed by the FMS; no truck can be dispatched to the dump point until the operator confirms the dozer has cleared and the new boundary is set.",
			},
		},
	})

	createSlide(database, s4.ID, "Why Capping is Essential", "CONTENT", 4, map[string]interface{}{
		"type":      "content",
		"heading":   "Why Capping is Essential",
		"introText": "Capping is the transition mechanism between OTE phases. Without it, the slope builds to its angle of repose and the face cannot advance further. The OTE-cap-dozer cycle is how a highwall dump progressively extends across its designed footprint.",
		"items": []map[string]interface{}{
			{
				"title":       "OTE Alone Cannot Advance the Face",
				"description": "OTE dumps build the slope from the existing crest. The slope eventually reaches its natural angle and further OTE tips just add height, not forward extent. Capping breaks this limit by bringing the berm forward so the next OTE sequence tips from a new, further-out position.",
			},
			{
				"title":       "The Repeating Cycle is the Dump's Lifespan",
				"description": "Each OTE-to-cap cycle advances the face by one increment. The dump's total capacity and operational life is determined by how many of these cycles it can complete before reaching the mine plan boundary. Efficient cycle management — minimising closure duration, maintaining count accuracy — directly affects dump throughput.",
			},
			{
				"title":       "Geotechnical Timing",
				"description": "Cap timing is also a geotechnical control. The number of OTE dumps before each cap is set by the geotechnical design — it controls face angle and slope stability. Changing the OTE count without geotechnical review changes the slope geometry and may compromise dump stability.",
			},
		},
	})

	// ── Section 5: Paddock Dumps ────────────────────────────────────────
	s5 := createSection(database, module.ID, "Paddock Dumps", 4)

	createSlide(database, s5.ID, "Paddock Dump from Above", "IMAGE", 0, map[string]interface{}{
		"type":        "image",
		"heading":     "Paddock Dump from Above",
		"description": "A paddock dump builds material in horizontal layers on flat or prepared ground. From the air, the completed lifts are visible as graduated steps behind the active tipping face — each layer represents one lift of material placed, compacted, and graded before the next lift begins.",
		"imageUrl":    "https://assets.miningacademy.com/images/paddock-dump-aerial-lifts.jpg",
	})
	createSlide(database, s5.ID, "What is a Paddock Dump?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What is a Paddock Dump?",
		"introText": "A paddock dump is any dump where material is placed on relatively flat or open ground and built up in controlled layers. Paddock dumping is one of the most versatile methods in mining — it is used for permanent waste disposal, ROM stockpile construction, building new access ramps, and forming elevated terraces or platforms on a mine site.",
		"items": []map[string]interface{}{
			{
				"title":       "Waste Dumps",
				"description": "Out-of-pit waste dumps are the most commonly pictured paddock dump application. Overburden is placed in controlled lifts and built into a permanent waste landform. Lift height and slope geometry are set by geotechnical design.",
			},
			{
				"title":       "ROM Stockpiles",
				"description": "Ore dumped onto a ROM pad ahead of the crusher is a paddock dump. Trucks tip to build and maintain the stockpile; a front-end loader or reclaim feeder draws material down for crusher feed. The same dump technique — controlled position, berm at the active face — applies.",
			},
			{
				"title":       "Ramps and Elevated Terraces",
				"description": "New haul road ramps, batter fills, and elevated working platforms are constructed by paddock dumping. A dozer or grader shapes the material after tipping. These applications are common during mine development and when mining geometry requires raising the floor of a working area.",
			},
		},
	})

	createSlide(database, s5.ID, "Paddock Dump Operations", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Paddock Dump Operations",
		"introText": "Operating on a paddock dump requires the same berm and approach discipline as any dump point, but with additional attention to lift boundaries and surface drainage. The active dump face moves progressively as each area of the current lift is filled.",
		"items": []map[string]interface{}{
			{
				"title":       "Following the Active Face",
				"description": "Dump trucks are dispatched to specific sections of the active lift face. As each section fills, the dispatch point moves to the next unfilled section. Operators must follow FMS instructions — tipping in the wrong section disrupts the planned lift sequence and drainage layout.",
			},
			{
				"title":       "Surface Grades and Drainage",
				"description": "The top surface of each lift must be graded to drain surface water away from the edge. Water ponding on a dump lift softens the surface and can trigger slumping. After heavy rain, dump surfaces must be inspected and re-graded before resuming.",
			},
			{
				"title":       "Berm Maintenance on Paddock Dumps",
				"description": "Paddock dump berms can be damaged by the push-over process — the dozer may inadvertently break berm continuity. After each push-over cycle, the berm must be rebuilt before the next truck reverses to that section.",
			},
		},
	})

	createSlide(database, s5.ID, "Paddock Dump Design Considerations", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Paddock Dump Design Considerations",
		"introText": "The design parameters for a paddock dump vary significantly with its purpose. A permanent waste landform has different requirements to a temporary ROM stockpile or a construction fill. What stays constant is the need for a correctly configured dump boundary and proper berm control at the active face.",
		"items": []map[string]interface{}{
			{
				"title":       "Lift Heights and Slope Geometry",
				"description": "Lift height and overall slope angle are geotechnical design outputs — they depend on material type, moisture, the purpose of the dump, and the site's closure requirements. Permanent waste landforms are designed to last beyond mine closure; temporary stockpiles and construction fills have different stability requirements. Always follow the site-specific design.",
			},
			{
				"title":       "Drainage",
				"description": "Surface drainage is critical regardless of dump purpose. Water ponding on a paddock dump softens material and can trigger surface failure. Permanent waste dumps require engineered catch drains and perimeter diversion drains. ROM stockpiles and temporary fills at minimum require a graded surface that sheds water away from the active face.",
			},
			{
				"title":       "Active Face Management",
				"description": "The active dump face on a paddock dump moves as each section fills. The FMS dispatch point must track the current face position. For AHS, the dump boundary polygon must be updated as the face advances — paddock dump faces can shift in ways that invalidate an existing boundary just as highwall dumps do.",
			},
		},
	})

	// ── Section 6: Crusher Dumps ────────────────────────────────────────
	s6 := createSection(database, module.ID, "Crusher Dumps", 5)

	createSlide(database, s6.ID, "Crusher Pocket", "IMAGE", 0, map[string]interface{}{
		"type":        "image",
		"heading":     "The Crusher Pocket",
		"description": "A crusher dump delivers ore directly to the primary crusher opening. Unlike open dump faces, the crusher pocket is a fixed concrete or steel structure with exact positional requirements. The truck must reverse to within centimetres of the correct position — there is no berm to stop against, only the backboard behind.",
		"imageUrl":    "https://assets.miningacademy.com/images/crusher-pocket-truck-dumping.jpg",
	})
	createSlide(database, s6.ID, "What is a Crusher Dump?", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "What is a Crusher Dump?",
		"introText": "A crusher dump is the delivery of run-of-mine (ROM) ore directly to the primary crusher pocket or ROM stockpile pad. Unlike waste dumps, crusher dumps deliver material to a fixed, precision location where downstream processing depends on consistent, uninterrupted feed.",
		"items": []map[string]interface{}{
			{
				"title":       "The Crusher Pocket",
				"description": "The crusher pocket is a hopper or grizzly opening at the top of the primary crusher. Trucks reverse to the edge of the pocket and tip directly in. The pocket has a limited buffer capacity — if trucks dump too fast, material backs up. If they dump too slowly, the crusher starves and throughput drops.",
			},
			{
				"title":       "ROM Pad Alternative",
				"description": "Where the crusher cannot accept direct truck feed (high-moisture ores, blending requirements, or crusher maintenance windows), ore is dumped onto a ROM stockpile pad adjacent to the crusher. A front-end loader or reclaim feeder then delivers ore to the crusher at a controlled rate.",
			},
			{
				"title":       "Fixed Infrastructure Hazards",
				"description": "Crusher dumps involve fixed infrastructure — concrete walls, backboards, haul road edges, and crusher pocket openings. The margin for positional error is much smaller than at open dump faces. GPS-assisted positioning and proximity detection are standard controls at most crusher dump points.",
			},
		},
	})

	createSlide(database, s6.ID, "Crusher Dump Operations", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Crusher Dump Operations",
		"introText": "Crusher dump approach and tipping technique are more constrained than open dump operations. The truck must position accurately, avoid overloading the pocket, and clear promptly so the next truck can access the dump point without delay.",
		"items": []map[string]interface{}{
			{
				"title":       "GPS-Guided Positioning",
				"description": "Most modern crusher dumps use GPS or proximity guidance systems to direct the truck to the correct reverse position. The system displays the truck position relative to the dump boundary and alerts the operator if they are outside the target zone. Manual positioning must be verified against physical reference marks.",
			},
			{
				"title":       "Controlled Tip Rate",
				"description": "When tipping into a crusher pocket, operators may need to pause mid-tip if the pocket is backing up — particularly with sticky or high-moisture ores. Tipping too fast into a backed-up pocket can cause material to rebound off the crusher grate or overflow the pocket walls.",
			},
			{
				"title":       "Clear and Clear Fast",
				"description": "After tipping, the truck must lower its body fully and move forward immediately to clear the dump bay for the next truck. Crusher dump points often have two-truck bays — but even then, queue management requires rapid turnaround to maintain crusher feed continuity.",
			},
		},
	})

	createSlide(database, s6.ID, "Crusher Dump Safety", "CONTENT", 3, map[string]interface{}{
		"type":      "content",
		"heading":   "Crusher Dump Safety",
		"introText": "Crusher dump safety is driven by proximity to fixed infrastructure and people. Crusher operators, lube technicians, and maintenance personnel may be working on or near the crusher structure at the same time as trucks are dumping. Isolation and communication protocols are essential.",
		"items": []map[string]interface{}{
			{
				"title":       "Crusher Isolation Before Dumping",
				"description": "No truck dumps into a crusher pocket while personnel are inside the crusher structure or on the dump platform. The crusher operation should have a clear entry isolation procedure — personnel working on the crusher must be signed out of the crusher zone before dump trucks access the bay.",
			},
			{
				"title":       "The Backboard and Spill Management",
				"description": "A backboard behind the crusher dump position prevents trucks from reversing too far and falling into the pocket. Large-fragment spills over the backboard must be cleaned up before they accumulate to a level that compromises the backboard's effectiveness. Spill material on the dump platform is a slip and trip hazard for personnel.",
				"imageUrl":    "https://assets.miningacademy.com/images/crusher-dump-backboard.jpg",
			},
			{
				"title":       "Oversize Material Hazard",
				"description": "Oversize rocks in the crusher feed can bridge the crusher pocket opening or rebound out of the pocket during dumping. When tipping oversize material, operators should tip slowly and observe for rebound before fully raising the body. Do not tip oversize material into an occupied crusher bay.",
			},
		},
	})

	// ── Section 7: Road Maintenance Dumps ──────────────────────────────
	s7 := createSection(database, module.ID, "Road Maintenance Dumps", 6)

	createSlide(database, s7.ID, "Road Maintenance Dump Operations", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "Road Maintenance Dump Operations",
		"introText": "Road maintenance dumps place construction or maintenance material — gravel, crushed rock, run-of-crusher fines — directly onto haul road surfaces for grading and compaction. This is a low-risk dump type when performed correctly, but the road environment creates unique hazards not present at fixed dump points.",
		"items": []map[string]interface{}{
			{
				"title":       "Tip on the Road, Not the Edge",
				"description": "Material must be placed in the centre of the road — not at the edge or over the berm. Edge-tipping destabilises the road shoulder, creates drop-off hazards, and makes the material harder for the grader to incorporate into the road surface.",
			},
			{
				"title":       "Traffic Management",
				"description": "Road maintenance tipping takes place on an active haul road. Other trucks and equipment are moving on the same surface. A traffic management plan or escort is required — tipping without traffic control puts both the operator and approaching traffic at risk. Flag or radio the approaching trucks before stopping to tip.",
			},
			{
				"title":       "Follow-Up Grading",
				"description": "Material tipped for road maintenance is not usable until graded and compacted. Lumps of raw material on the haul road surface are a tyre damage hazard. Graders should follow the tipping truck within the same shift to incorporate the material before the next production fleet cycle.",
			},
		},
	})

	createSlide(database, s7.ID, "Material Selection for Road Maintenance", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Material Selection for Road Maintenance",
		"introText": "The wrong material on a haul road is worse than no material. Using highly weathered, clay-rich, or poorly graded material for road maintenance can create muddy, unstable surfaces that damage tyres, reduce rolling resistance, and increase incident risk in wet conditions.",
		"items": []map[string]interface{}{
			{
				"title":       "Wearing Course Material",
				"description": "The haul road wearing course requires well-graded crushed rock or natural gravel with a high percentage of fines (to bind) and coarse aggregate (for strength). Material specifications are defined in the haul road design — typically a CBR > 80 and plasticity index < 12 for the wearing course.",
			},
			{
				"title":       "Subgrade Repair Material",
				"description": "Soft spots and potholes in the subgrade need to be cleaned out and backfilled with suitable material before surface resealing. Tipping surface material onto a soft subgrade just covers the symptom — the pothole will return within a shift. Report subgrade failures to the road maintenance supervisor.",
			},
			{
				"title":       "Segregation Awareness",
				"description": "Large fragments segregate during hauling and tipping — coarse material rolls to the outside, fines concentrate in the middle. After tipping, the grader must blend the material from both edges of the windrow to re-achieve the correct gradation before compaction.",
			},
		},
	})

	// ── Section 8: Autonomous Dumping ───────────────────────────────────
	s8 := createSection(database, module.ID, "Autonomous Dumping", 7)

	createSlide(database, s8.ID, "AHS and the Dump Workflow", "CONTENT", 0, map[string]interface{}{
		"type":      "content",
		"heading":   "AHS and the Dump Workflow",
		"introText": "In an autonomous haulage system, the dump workflow is the same as in manned operations — but the truck has no operator. Every decision a manned driver makes at the dump point (where to stop, when to tip, when to leave) is encoded in configuration and executed by the system. The FMS orchestrates the full lifecycle of each dump point.",
		"items": []map[string]interface{}{
			{
				"title":       "The Dump Point as a Configured Object",
				"description": "Each dump point in an AHS is a configurable entity in the FMS. It has a boundary (where trucks are allowed), one or more scripts (what the truck does when it arrives), and state (open, closed, capping). The FMS dispatches trucks to dump points and enforces the configured constraints.",
			},
			{
				"title":       "Scripts Define the Truck's Behaviour",
				"description": "A dump script defines the truck's reverse path, tip position, body sequence, and exit path. A highwall dump point has at least two scripts: an OTE dump script (tip over the edge) and a capping dump script (tip at the crest). The FMS selects which script to use based on the dump point's current state.",
			},
			{
				"title":       "Human Roles Shift to Configuration and Authorisation",
				"description": "In AHS, the dump controller and FMS operator take on the work that operators previously did at the face. They configure dump points, manage state transitions, authorise boundary changes, and decide when the dump is safe to reopen. Every decision that was made in the cab is now made at the console — before the truck arrives.",
			},
		},
	})

	createSlide(database, s8.ID, "Highwall Dump State Management", "CONTENT", 1, map[string]interface{}{
		"type":      "content",
		"heading":   "Highwall Dump State Management",
		"introText": "A highwall dump point in an AHS passes through a defined sequence of states during each OTE-to-cap cycle. The FMS manages these transitions automatically where it can — and requires human authorisation where it cannot. Getting this state machine right is the core challenge of AHS highwall dump integration.",
		"items": []map[string]interface{}{
			{
				"title":       "OTE State — Counting Dumps",
				"description": "The dump point is open under the OTE script. The FMS dispatches trucks and increments the dump count with each tip. When the count reaches the configured threshold, the FMS automatically transitions the point to capping state — no human action needed for this transition.",
			},
			{
				"title":       "Capping State — One Cap Per Point",
				"description": "The dump point switches to the capping dump script. The next truck dispatched tips its load at the crest rather than over the edge. Once the capping dump is complete at this point, the point moves to a waiting state. When all points in the face have capped, the FMS transitions the entire face to closed.",
			},
			{
				"title":       "Closed State — Dozer Work",
				"description": "No trucks are dispatched to any point in the face. The dump controller assigns the dozer task. The dozer works across all points, pushing the berm outward. Once the dozer operator confirms completion, the dump controller updates the boundary and reopens the dump — resetting all points to OTE state with a zero count.",
			},
		},
		"keyFact": "The closed state is not a failure mode — it is a planned phase of the advance cycle. Designing the fleet schedule to absorb the closure window is an operational planning task, not a system limitation to be designed around.",
	})

	createSlide(database, s8.ID, "Configuring a Highwall Dump Point", "CONTENT", 2, map[string]interface{}{
		"type":      "content",
		"heading":   "Configuring a Highwall Dump Point",
		"introText": "When setting up a new highwall dump point in the FMS, the operator must configure both the OTE and capping aspects of that point before any truck is dispatched. The configuration defines both the physical parameters (boundary, position) and the operational logic (which script, what count).",
		"items": []map[string]interface{}{
			{
				"title":       "OTE Dump Script",
				"description": "The OTE script defines the reverse path to the berm, the tip position relative to the boundary, and the body raise/lower sequence. It also captures the dump count threshold — the number of OTE tips before this point transitions to capping. This is a geotechnically-derived number that must not be changed without review.",
			},
			{
				"title":       "Capping Dump Script",
				"description": "The capping script defines a different reverse path and tip position from the OTE script — the truck tips at the crest rather than at the berm edge. The positional difference between the two scripts is what distinguishes a capping dump from an OTE dump. Both scripts reference the same dump point boundary.",
			},
			{
				"title":       "The Boundary Polygon",
				"description": "The boundary polygon defines the zone within which the dump point's scripts are valid. After each dozer advance cycle, this polygon must be updated to reflect the new berm position before the dump reopens. Updating the boundary is a human-authorised step — the FMS will not automatically infer a new boundary from the dozer's path.",
			},
		},
	})

	createSlide(database, s8.ID, "Manned vs Autonomous Dumping", "COMPARISON", 3, map[string]interface{}{
		"type":    "comparison",
		"heading": "Manned vs Autonomous Dumping",
		"items": []map[string]interface{}{
			{
				"name":        "Manned Dumping",
				"description": "Human operator in the cab executes the dump procedure with direct visual and sensory feedback.",
				"pros": []string{
					"Operator adapts in real time to conditions — berm damage, soft ground, poor visibility",
					"No configuration latency — driver adjusts position on the fly",
					"Natural situational awareness at the face without sensor dependency",
					"Simpler system — no FMS state machine or script management required",
				},
				"cons": []string{
					"Human fatigue and attention variability increase error rate on long shifts",
					"No automated record of dump position or count per cycle",
					"Operators must be present at high-risk dump faces",
					"OTE-to-cap transitions rely on supervisor instruction rather than system logic",
				},
			},
			{
				"name":        "Autonomous Dumping (AHS)",
				"description": "FMS-managed truck executes the dump based on configured scripts with no operator in the cab.",
				"pros": []string{
					"Consistent execution every cycle — no fatigue or attention effects",
					"Full audit log of every dump event, count, and state transition",
					"Removes personnel from high-risk dump environments",
					"OTE-to-cap transition is automatic once count threshold is reached",
				},
				"cons": []string{
					"Requires precise, maintained configuration before any truck is dispatched",
					"Dump closure for dozer work must be planned into fleet scheduling",
					"Boundary updates require human authorisation — creates a gap between physical advance and FMS state",
					"System cannot adapt to undetected physical changes without sensor or human input",
				},
			},
		},
	})

	// ── Section 9: Assessment ───────────────────────────────────────────
	s9 := createSection(database, module.ID, "Assessment", 8)

	quizDesc := "Test your understanding of dump operations — standard procedures, highwall and capping mechanics, dump types, and autonomous dumping considerations."
	quiz := models.Quiz{
		ID:               uuid.New().String(),
		Title:            "Dump Operations Assessment",
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
			question:    "A haul truck arrives at a highwall dump point and finds no berm present. What is the correct action?",
			qType:       models.QuestionMultipleChoice,
			explanation: "No berm means no dump — the primary engineered safeguard is absent. The operator must park clear of the dump area and report to the supervisor before the next truck approaches. Continuing to tip without a berm removes the only physical barrier preventing the truck from going over the edge.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Tip carefully at reduced speed", false},
				{"Park clear and report to the supervisor immediately", true},
				{"Use the previous tyre tracks as a position guide", false},
				{"Call the spotter to stand at the edge instead", false},
			},
		},
		{
			question:    "What is the minimum berm height required at a haul truck dump point?",
			qType:       models.QuestionMultipleChoice,
			explanation: "The standard minimum berm height is half the inflated diameter of the largest rear tyre using that dump point. For a large haul truck with 4.0 m tyres, this means a minimum 2.0 m berm. A lower berm can be overridden by the weight of the loaded truck body.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"One full tyre diameter", false},
				{"Half the tyre diameter", true},
				{"1.0 m regardless of tyre size", false},
				{"Any height is acceptable if a spotter is present", false},
			},
		},
		{
			question:    "What is the primary purpose of a capping dump at a highwall?",
			qType:       models.QuestionMultipleChoice,
			explanation: "The capping dump delivers material to the crest for the dozer to push outward. The dozer reconstructs the berm at the new advanced position — the capping dump provides the material; the dozer advances the face.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"To close out the dump point permanently", false},
				{"To deliver material at the crest for the dozer to push the berm outward", true},
				{"To drive onto the toe and raise the dump platform", false},
				{"To compact the dump surface for vehicle access", false},
			},
		},
		{
			question:    "After all dump points in a highwall face have completed their capping dumps, what happens next?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Once every point in the face has received its capping dump, the FMS closes the dump entirely. A dozer then works across all points, pushing the capping material outward to advance and reconstruct the berm. No trucks are dispatched during this phase.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Trucks resume OTE immediately from the new position", false},
				{"The FMS closes the dump and a dozer reconstructs the berm across all points", true},
				{"Survey updates the boundary and trucks continue without interruption", false},
				{"One truck per point performs a second capping dump", false},
			},
		},
		{
			question:    "In an AHS, a highwall dump point is configured with two scripts. What are they?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Each highwall dump point needs an OTE dump script (tip over the edge) and a capping dump script (tip at the crest for the dozer). The FMS selects which script to apply based on the dump point's current count state.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"A manned script and an autonomous script", false},
				{"An OTE dump script and a capping dump script", true},
				{"A forward dump script and a reverse dump script", false},
				{"A primary dump script and a backup dump script", false},
			},
		},
		{
			question:    "In an AHS highwall dump, the dozer has finished pushing the berm outward across all points. What must happen before trucks can resume OTE dumping?",
			qType:       models.QuestionMultipleChoice,
			explanation: "The dump boundary polygon in the FMS must be updated to reflect the new berm position, and the dump must be reopened by the controller with counts reset to zero. The FMS will not automatically infer a new boundary — it requires human authorisation.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Nothing — the FMS automatically detects the new berm and dispatches trucks", false},
				{"The boundary polygon must be updated and the controller must reopen the dump with counts reset", true},
				{"The first truck performs a slow reverse to confirm the new berm position", false},
				{"A spotter must stand at each point until three trucks have dumped safely", false},
			},
		},
		{
			question:    "A truck is reversing at a highwall dump with a spotter in place. The spotter's radio goes silent mid-reverse. What should the driver do?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Loss of spotter contact during a reverse is a stop condition. The driver must stop immediately and not resume the reverse until contact is re-established. A silent radio means the driver has lost the last line of defence at the dump edge.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Continue reversing slowly and check for the spotter visually", false},
				{"Stop immediately and do not move until contact is re-established", true},
				{"Complete the current dump and then report the radio issue", false},
				{"Reverse only as far as the previous truck's tyre tracks", false},
			},
		},
		{
			question:    "A new waste dump area is to be established on flat open ground 800 m from the active pit. There is no highwall edge at this location. What dump type should be used?",
			qType:       models.QuestionMultipleChoice,
			explanation: "Flat open ground without a highwall edge calls for a paddock dump — layered horizontal lifts constructed from the ground up. A highwall dump requires an active highwall edge for OTE placement. A crusher dump is for ROM ore. A road maintenance dump is for maintenance material only.",
			options: []struct {
				text      string
				isCorrect bool
			}{
				{"Highwall dump — it has the highest capacity", false},
				{"Crusher dump — it is the closest infrastructure", false},
				{"Paddock dump — flat ground layered lift construction", true},
				{"Road maintenance dump — material will eventually be used for roads", false},
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

	createSlide(database, s9.ID, "Dump Operations Assessment", "QUIZ", 0, map[string]interface{}{
		"type":    "quiz",
		"heading": "Dump Operations Assessment",
		"quizId":  quiz.ID,
		"inlineQuestions": []map[string]interface{}{
			{
				"question":     "What must be present before any haul truck tips at a dump point?",
				"options":      []string{"A spotter", "A berm of at least half the tyre diameter", "Radio contact with dispatch", "GPS confirmation"},
				"correctIndex": 1,
				"explanation":  "The berm is the primary engineered safeguard — minimum half the inflated tyre diameter. No berm, no dump.",
			},
			{
				"question":     "What does the capping dump provide, and who uses it?",
				"options":      []string{"A platform for the next OTE truck", "Material at the crest for the dozer to push the berm outward", "Compaction for the OTE slope", "A boundary marker for the FMS"},
				"correctIndex": 1,
				"explanation":  "The capping dump tips material at the crest — not over the edge. The dozer then pushes that material outward to advance and reconstruct the berm.",
			},
			{
				"question":     "In AHS, what triggers the automatic transition from OTE script to capping script at a dump point?",
				"options":      []string{"A supervisor instruction", "The FMS dump count reaching the configured threshold", "The truck's proximity sensor detecting the toe", "A survey update to the boundary polygon"},
				"correctIndex": 1,
				"explanation":  "The FMS tracks the dump count per point. When the count reaches the configured OTE threshold, the system automatically switches to the capping script for that point.",
			},
		},
	})

	// ── Section 10: Completion ──────────────────────────────────────────
	s10 := createSection(database, module.ID, "Completion", 9)

	createSlide(database, s10.ID, "Module Complete", "COMPLETION", 0, map[string]interface{}{
		"type":    "completion",
		"heading": "Module Complete!",
		"message": "You have completed the Dump Operations module.",
		"summary": "You can now describe the standard dump procedure and its mandatory controls, explain the OTE-to-cap cycle and why capping is essential for advancing a highwall dump, identify the correct dump type for a given material and site geometry, and explain how AHS changes dump boundary management responsibilities.",
	})

	log.Println("Dump Operations module seeded successfully.")
}
