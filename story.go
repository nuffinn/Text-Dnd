package main

type StoryEngine struct {
	Scenes       map[int]*Scene
	CurrentScene int
}

type Scene struct {
	ID          int
	Text        string
	Choices     []Choice
	IsEnding    bool
	EndingType  string
	NextSuccess int
	NextFailure int
}

func (s *Scene) GetNextScene(choiceID string, success bool) int {
	// Find the choice and determine next scene
	for _, choice := range s.Choices {
		if choice.ID == choiceID {
			if success {
				return s.NextSuccess
			}
			return s.NextFailure
		}
	}
	return s.NextSuccess
}

func NewStoryEngine() *StoryEngine {
	engine := &StoryEngine{
		Scenes:       make(map[int]*Scene),
		CurrentScene: 0,
	}
	engine.initializeStory()
	return engine
}

func (e *StoryEngine) GetScene(id int) *Scene {
	if scene, exists := e.Scenes[id]; exists {
		return scene
	}
	return e.Scenes[0]
}

func (e *StoryEngine) GetCurrentScene() *Scene {
	return e.GetScene(e.CurrentScene)
}

func (e *StoryEngine) initializeStory() {
	// Scene 0: Mission Start
	e.Scenes[0] = &Scene{
		ID: 0,
		Text: `🚀 MISSION LOG - ARTEMIS EXPEDITION

Year 2045. You are crew aboard the ARTEMIS, humanity's first commercial asteroid mining vessel. Your target: asteroid 16 Psyche, a metallic asteroid worth an estimated $10,000 quadrillion in rare metals.

After 3 years of travel, you've reached the asteroid. But as you approach, your sensors detect something unexpected...

⚠️ ALERT: Unknown metallic structure detected on asteroid surface. Not matching any known natural formation.

Captain's Log: "This changes everything. Command wants us to investigate before beginning mining operations."

The crew gathers on the bridge. What do you do?`,
		Choices: []Choice{
			{
				ID:          "investigate",
				Text:        "Send a reconnaissance drone to investigate the structure",
				SkillCheck:  "piloting",
				Difficulty:  12,
				Description: "Use the ship's drone for a closer look. Requires piloting check.",
			},
			{
				ID:          "scan",
				Text:        "Perform detailed sensor scans from safe distance",
				SkillCheck:  "science",
				Difficulty:  10,
				Description: "Analyze the structure with ship sensors. Requires science check.",
			},
			{
				ID:          "eva",
				Text:        "Prepare an EVA team to investigate directly",
				SkillCheck:  "engineering",
				Difficulty:  14,
				Description: "Suit up and go out there. Requires engineering check to prepare.",
			},
			{
				ID:          "report",
				Text:        "Report to Earth and await instructions (3-year round trip delay)",
				SkillCheck:  "",
				Difficulty:  0,
				Description: "Play it safe and wait for orders.",
			},
		},
		NextSuccess: 1,
		NextFailure: 2,
	}

	// Scene 1: Successful Investigation
	e.Scenes[1] = &Scene{
		ID: 1,
		Text: `📡 DISCOVERY

Your approach is careful and professional. As your instruments get clearer readings, the crew holds their breath...

The structure is artificial. Ancient, but artificial. It's a probe - clearly not human in origin. The design is elegant, mathematical, clearly built by intelligent beings. Radio-isotope dating suggests it's been here for approximately 50,000 years.

Dr. Chen, your scientist: "This is First Contact... or rather, first evidence. We need to be careful."

The probe suddenly activates, projecting a holographic star map with a trajectory leading to the outer solar system. Then it goes dark again.

⚠️ SHIP STATUS: Fuel at 60%, Oxygen at 85%

What is your next move?`,
		Choices: []Choice{
			{
				ID:          "extract",
				Text:        "Carefully extract the probe and bring it aboard for study",
				SkillCheck:  "engineering",
				Difficulty:  15,
				Description: "Recover the artifact. High risk, high reward.",
			},
			{
				ID:          "document",
				Text:        "Document everything and continue with mining mission",
				SkillCheck:  "science",
				Difficulty:  10,
				Description: "Gather data, report it, stick to the plan.",
			},
			{
				ID:          "follow",
				Text:        "Follow the trajectory shown in the star map",
				SkillCheck:  "piloting",
				Difficulty:  16,
				Description: "Change mission parameters. Very risky.",
			},
		},
		NextSuccess: 3,
		NextFailure: 4,
	}

	// Scene 2: Failed Initial Approach
	e.Scenes[2] = &Scene{
		ID: 2,
		Text: `⚠️ SYSTEMS FAILURE

Your approach goes wrong. The drone crashes into the asteroid surface, or the sensors malfunction, or the EVA team's tether snaps...

Emergency protocols activate. The crew scrambles to prevent disaster.

Engineer Torres: "We've lost 20% of our fuel reserves in the recovery operation!"

You manage to salvage the situation, but you've captured only partial data on the structure. It appears artificial, possibly a probe or beacon. But you've lost time and resources.

⚠️ SHIP STATUS: Fuel at 40%, Oxygen at 75%, Ship Integrity at 90%

The crew is shaken but intact. You need to make a decision fast.`,
		Choices: []Choice{
			{
				ID:          "mine",
				Text:        "Abandon investigation, begin mining operations immediately",
				SkillCheck:  "leadership",
				Difficulty:  12,
				Description: "Focus on the mission. Keep the crew focused.",
			},
			{
				ID:          "retry",
				Text:        "Attempt another investigation with remaining resources",
				SkillCheck:  "engineering",
				Difficulty:  16,
				Description: "High risk - you're low on resources.",
			},
			{
				ID:          "return",
				Text:        "Abort mission and return to Earth with findings",
				SkillCheck:  "",
				Difficulty:  0,
				Description: "The safe choice, but mission failure.",
			},
		},
		NextSuccess: 5,
		NextFailure: 6,
	}

	// Scene 3: Probe Recovered - Success Branch
	e.Scenes[3] = &Scene{
		ID: 3,
		Text: `🎯 BREAKTHROUGH

The probe is safely aboard. Your team works carefully, documenting every detail. The technology is far beyond current human capabilities, yet somehow comprehensible.

Dr. Chen makes a startling discovery: "The probe contains a data storage medium. And it's... it's accessible. Like it was meant to be found."

The data shows a warning: A massive coronal mass ejection from the Sun will occur in approximately 72 hours, directly in your current orbital path. The ancient builders left this warning system in place.

The probe also contains technical schematics - propulsion technology that could revolutionize space travel.

Navigator Kim: "If the warning is accurate, we need to move NOW. But if we take time to copy all the data..."

⚠️ CRITICAL DECISION - TIME SENSITIVE`,
		Choices: []Choice{
			{
				ID:          "evacuate",
				Text:        "Evacuate immediately, take probe with minimal data copied",
				SkillCheck:  "piloting",
				Difficulty:  13,
				Description: "Save the crew and hardware, lose most of the data.",
			},
			{
				ID:          "copy",
				Text:        "Take time to copy all data, then escape",
				SkillCheck:  "science",
				Difficulty:  17,
				Description: "Risky - you might not make it. But the data could save humanity.",
			},
			{
				ID:          "shield",
				Text:        "Use the probe's technology to shield the ship",
				SkillCheck:  "engineering",
				Difficulty:  18,
				Description: "Extremely risky - attempt to activate alien technology.",
			},
		},
		NextSuccess: 7,
		NextFailure: 8,
	}

	// Scene 4: Probe Study Failed
	e.Scenes[4] = &Scene{
		ID: 4,
		Text: `⚡ CLOSE CALL

The extraction goes wrong. The probe's protective systems activate, releasing an electromagnetic pulse that temporarily disables your ship's systems.

Emergency power kicks in after 3 terrifying minutes. Everyone is safe, but the probe is damaged and the ship's navigation systems need recalibration.

Chief Engineer Torres: "We're lucky to be alive. That thing was defending itself."

While repairing systems, your communication officer picks up unusual solar activity. A massive coronal mass ejection is building - and it's aimed right at your position.

⚠️ SHIP STATUS: Fuel at 50%, Systems at 60%, Time to CME impact: 48 hours

You have limited options with damaged systems.`,
		Choices: []Choice{
			{
				ID:          "emergency",
				Text:        "Emergency burn to escape blast radius",
				SkillCheck:  "piloting",
				Difficulty:  15,
				Description: "Risky with damaged systems, but fastest option.",
			},
			{
				ID:          "repair",
				Text:        "Repair systems first, then controlled evacuation",
				SkillCheck:  "engineering",
				Difficulty:  16,
				Description: "Safer but time-consuming.",
			},
			{
				ID:          "shelter",
				Text:        "Take shelter in asteroid's shadow",
				SkillCheck:  "science",
				Difficulty:  14,
				Description: "Use the asteroid as a shield.",
			},
		},
		NextSuccess: 9,
		NextFailure: 10,
	}

	// Scene 5: Mining Success After Recovery
	e.Scenes[5] = &Scene{
		ID: 5,
		Text: `⛏️ OPERATION SUCCESSFUL

Your crew pulls together brilliantly. Despite the earlier setback, you refocus on the primary mission. The mining operation proceeds smoothly.

Over the next weeks, you extract valuable platinum-group metals worth billions. The alien structure is documented but left in place - above your pay grade to handle.

As you prepare for the journey home, your sensors detect unusual activity from the structure. It's transmitting a signal... directly at your ship.

The message is simple, mathematical - a sequence of prime numbers, then coordinates. They point to Neptune's orbit.

Captain's decision time: You've completed the mission successfully. But this signal...

⚠️ SHIP STATUS: Fuel at 55%, Cargo hold full of valuable metals, All systems nominal`,
		Choices: []Choice{
			{
				ID:          "home",
				Text:        "Return to Earth with the mining haul and data",
				SkillCheck:  "",
				Difficulty:  0,
				Description: "Mission complete. Take the win.",
			},
			{
				ID:          "investigate_neptune",
				Text:        "Investigate the coordinates near Neptune",
				SkillCheck:  "leadership",
				Difficulty:  15,
				Description: "Risky, but potentially historic.",
			},
			{
				ID:          "beacon",
				Text:        "Leave a beacon and return to Earth for a proper expedition",
				SkillCheck:  "science",
				Difficulty:  10,
				Description: "Responsible middle ground.",
			},
		},
		NextSuccess: 11,
		NextFailure: 12,
	}

	// Scene 6: Critical Failure - Resource Crisis
	e.Scenes[6] = &Scene{
		ID: 6,
		Text: `🔴 CRISIS

Everything that could go wrong did. The second attempt fails catastrophically. You've now lost critical resources, and crew morale is at an all-time low.

⚠️ SHIP STATUS: Fuel at 25%, Oxygen at 60%, Systems at 70%

Medical Officer Patel: "We have enough resources to make it back to Earth, but only if we leave NOW. Any delay and we're dead in space."

But you've just detected that the alien structure is actually a warning beacon - and it's alerting you to something. Solar sensors show a massive coronal mass ejection heading your way. ETA: 36 hours.

This is it. The moment that defines whether you all make it home.`,
		Choices: []Choice{
			{
				ID:          "immediate",
				Text:        "Immediate emergency departure for Earth",
				SkillCheck:  "piloting",
				Difficulty:  17,
				Description: "Burn everything you have. It's now or never.",
			},
			{
				ID:          "distress",
				Text:        "Send distress signal and hope for rescue",
				SkillCheck:  "leadership",
				Difficulty:  18,
				Description: "Nearest ship is months away. Long shot.",
			},
			{
				ID:          "bunker",
				Text:        "Create emergency shelter in the asteroid",
				SkillCheck:  "engineering",
				Difficulty:  19,
				Description: "Desperate measure. Mine into the asteroid for protection.",
			},
		},
		NextSuccess: 13,
		NextFailure: 14,
	}

	// Scene 7: Best Ending - Data Secured
	e.Scenes[7] = &Scene{
		ID:   7,
		Text: `🌟 TRIUMPH OF HUMANITY

Against all odds, you succeed. Your team works with precision and determination, copying the alien data while preparing for emergency departure.

The escape is dramatic - you leave the asteroid mere hours before the coronal mass ejection hits. But you make it.

The probe's data is revolutionary. Within a decade, humanity has new propulsion systems, renewable energy sources, and a roadmap to the stars. The ancient builders left a gift, and you ensured humanity received it.

Your crew becomes legendary. The Artemis Expedition is taught in schools. First Contact wasn't with living aliens, but with their legacy - and it changed everything.

Years later, a follow-up mission to Neptune's coordinates finds something extraordinary: a library, left by the same builders, containing knowledge from a thousand worlds.

🏆 ENDING: "The Seekers' Legacy" - BEST ENDING
✨ Mission Success: 100%
👥 All Crew Survived
🎖️ Crew Achievement: Heroes of Humanity

The stars are calling. And now, humanity can answer.`,
		IsEnding:   true,
		EndingType: "Victory - Humanity reaches for the stars",
	}

	// Scene 8: Bittersweet Ending - Crew Saved, Data Lost
	e.Scenes[8] = &Scene{
		ID: 8,
		Text: `⭐ SURVIVORS

You make the hard call: crew first, data second. The probe comes with you, but most of the data is lost in the hasty departure.

The escape is harrowing. The coronal mass ejection hits just as you clear the danger zone. One of your crew members is injured by radiation exposure but survives thanks to quick medical action.

The probe itself is studied for decades. While much was lost, enough remains to advance human technology by years. Your testimony about the star map and the warning system leads to new deep-space monitoring programs.

You're celebrated as heroes, not for what you brought back, but for bringing everyone home alive. Sometimes the greatest discovery is rediscovering the value of human life.

The asteroid keeps its secrets. But you kept your crew.

🏆 ENDING: "Those Who Return" - GOOD ENDING
✨ Mission Success: 75%
👥 All Crew Survived
🎖️ Crew Achievement: Made the Hard Choices`,
		IsEnding:   true,
		EndingType: "Good Ending - You brought everyone home",
	}

	// Scene 9: Narrow Escape
	e.Scenes[9] = &Scene{
		ID: 9,
		Text: `🎯 ESCAPE VELOCITY

Your quick thinking and skilled crew pull off the impossible. Whether through a desperate burn, smart repairs, or using the asteroid itself as a shield, you escape the solar storm.

The return journey is long and quiet. You've survived, but with limited data on the alien probe and no mining haul. The mission is a technical failure but a human triumph.

Earth welcomes you as survivors of an incredible journey. Your report on the alien artifact leads to future missions, better prepared and better equipped.

Ten years later, you watch as the new expedition launches toward the asteroid. Your data made it possible. You didn't finish the mission, but you started something greater.

🏆 ENDING: "Against the Odds" - SURVIVAL ENDING
✨ Mission Success: 50%
👥 All Crew Survived
🎖️ Crew Achievement: Survivors`,
		IsEnding:   true,
		EndingType: "Survival Ending - You lived to tell the tale",
	}

	// Scene 10: Tragic Ending
	e.Scenes[10] = &Scene{
		ID: 10,
		Text: `⚫ SILENT STARS

Despite your best efforts, you can't escape in time. The coronal mass ejection hits the Artemis with full force.

Systems fail. Life support struggles. The crew fights for survival, but in the depths of space, sometimes fighting isn't enough.

Your final transmission reaches Earth three years later:

"This is Captain [Name] of the Artemis. If you're receiving this, we didn't make it. But we found something incredible - proof we're not alone. The coordinates are in our database. Don't let our sacrifice be in vain. Reach for the stars. Someone out there left us a message. Make sure we answer it..."

The mission is declared lost. But your sacrifice isn't forgotten. The Artemis II expedition launches five years later, carrying your names and continuing your work.

In the end, you became part of the journey to the stars.

🏆 ENDING: "Stardust" - TRAGIC ENDING
✨ Mission Success: 25%
💀 Crew Status: Lost
🎖️ Crew Achievement: Martyrs of Exploration

In space, everyone can hear you scream. But they can also hear you dream.`,
		IsEnding:   true,
		EndingType: "Tragic Ending - The price of exploration",
	}

	// Scenes 11-14: Additional endings
	e.Scenes[11] = &Scene{
		ID: 11,
		Text: `🌌 FIRST CONTACT

You make the bold decision to investigate Neptune's coordinates. The journey takes months, burning through resources. The crew debates whether this is brilliance or madness.

As you approach the coordinates, your sensors detect something massive: a structure orbiting Neptune, ancient and vast. It's a library - a repository of knowledge from dozens of extinct civilizations.

The builders left this for whoever would find it. A gift to young species reaching for the stars.

Your discovery reshapes human civilization. This is First Contact with the legacy of a thousand worlds.

🏆 ENDING: "The Library of Worlds" - LEGENDARY ENDING
✨ Mission Success: 150%
👥 All Crew Survived
🎖️ Crew Achievement: Legends of the Cosmos`,
		IsEnding:   true,
		EndingType: "Legendary Ending - You found something greater",
	}

	e.Scenes[12] = &Scene{
		ID: 12,
		Text: `📍 THE LONG GAME

You make the smart call: beacon deployed, mission complete, return home to prepare properly.

The journey home takes three years. You return as heroes with both alien data and mining success. The beacon transmits coordinates back to Earth.

Fifteen years later, you're selected to captain the Artemis II - a proper expedition to Neptune's coordinates, equipped with everything you'll need.

As you embark on this new journey, you reflect: sometimes the greatest discovery is knowing there's more to discover.

🏆 ENDING: "The Patient Explorer" - OPTIMAL ENDING
✨ Mission Success: 120%
👥 All Crew Survived
🎖️ Crew Achievement: Wisdom and Courage`,
		IsEnding:   true,
		EndingType: "Optimal Ending - The smart play",
	}

	e.Scenes[13] = &Scene{
		ID: 13,
		Text: `🚀 AGAINST ALL ODDS

Battered, exhausted, running on fumes and prayers - you make it. The Artemis limps back to Earth orbit after a harrowing journey.

Two crew members need extended medical care. The ship will never fly again. But you all survived.

The alien artifact data was lost. The mining haul was abandoned. On paper, the mission failed.

But you brought everyone home.

Earth celebrates not your discovery, but your determination. You prove that humanity's greatest strength isn't technology - it's the refusal to give up on each other.

🏆 ENDING: "No One Left Behind" - HEROIC ENDING
✨ Mission Success: 60%
👥 All Crew Survived
🎖️ Crew Achievement: True Heroes`,
		IsEnding:   true,
		EndingType: "Heroic Ending - The crew is what matters",
	}

	e.Scenes[14] = &Scene{
		ID: 14,
		Text: `⚫ THE PRICE OF STARS

The final log entry from the Artemis is found years later by a deep-space probe:

"Day 847. Oxygen at 2%. We tried everything. The shelter plan failed. We're too far for rescue. Too damaged to run.

But we don't regret trying. We saw wonders. We found proof of other minds in the cosmos. We pushed the boundary of human achievement.

If someone finds this: we were the Artemis crew. We reached for the stars. Some of us fell. But the reaching matters.

Keep reaching."

Your sacrifice leads to better ships, better training, better protocols. Future crews survive because you didn't.

In the vacuum of space, your legacy echoes forever.

🏆 ENDING: "The Reaching" - MEMORIAL ENDING
✨ Mission Success: 0%
💀 Crew Status: Lost
🎖️ Crew Achievement: Pioneers

Ad astra per aspera - Through hardship to the stars.`,
		IsEnding:   true,
		EndingType: "Memorial Ending - Some prices are too high",
	}
}
