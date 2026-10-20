-- Additive installation: run after the base schema, with the server stopped.
-- Definitions are shared catalogue data; player card instances remain in match snapshots.
-- Re-running preserves existing rows and never overwrites catalogue content.
begin immediate;

create table if not exists card_designs (
    catalogue_version text not null check (length(trim(catalogue_version)) > 0),
    design_id text not null check (length(trim(design_id)) > 0),
    rules_version text not null check (length(trim(rules_version)) > 0),
    sort_order integer not null check (typeof(sort_order) = 'integer' and sort_order >= 0),
    pool text not null default '',
    name text not null default '' check (length(trim(name)) > 0),
    text text not null default '',
    timing text not null default '',
    duration integer not null default 0 check (typeof(duration) = 'integer' and duration >= 0),
    element text not null default '',
    victory_path text not null default '',
    persistence text not null default '',
    disposal text not null default '',
    placeholder integer not null default 0 check (placeholder in (0, 1)),
    cost_wis integer check (cost_wis is null or (typeof(cost_wis) = 'integer' and cost_wis >= 0)),
    cost_wealth integer check (cost_wealth is null or (typeof(cost_wealth) = 'integer' and cost_wealth >= 0)),
    play_requirement_fire integer check (play_requirement_fire is null or (typeof(play_requirement_fire) = 'integer' and play_requirement_fire >= 0)),
    play_requirement_earth integer check (play_requirement_earth is null or (typeof(play_requirement_earth) = 'integer' and play_requirement_earth >= 0)),
    play_requirement_air integer check (play_requirement_air is null or (typeof(play_requirement_air) = 'integer' and play_requirement_air >= 0)),
    play_requirement_water integer check (play_requirement_water is null or (typeof(play_requirement_water) = 'integer' and play_requirement_water >= 0)),
    play_requirement_nobility integer check (play_requirement_nobility is null or (typeof(play_requirement_nobility) = 'integer' and play_requirement_nobility >= 0)),
    play_requirement_church integer check (play_requirement_church is null or (typeof(play_requirement_church) = 'integer' and play_requirement_church >= 0)),
    play_requirement_infernal integer check (play_requirement_infernal is null or (typeof(play_requirement_infernal) = 'integer' and play_requirement_infernal >= 0)),
    play_requirement_might integer check (play_requirement_might is null or (typeof(play_requirement_might) = 'integer' and play_requirement_might >= 0)),
    play_requirement_popularity integer check (play_requirement_popularity is null or (typeof(play_requirement_popularity) = 'integer' and play_requirement_popularity >= 0)),
    expertise_fire integer check (expertise_fire is null or (typeof(expertise_fire) = 'integer' and expertise_fire >= 0)),
    expertise_earth integer check (expertise_earth is null or (typeof(expertise_earth) = 'integer' and expertise_earth >= 0)),
    expertise_air integer check (expertise_air is null or (typeof(expertise_air) = 'integer' and expertise_air >= 0)),
    expertise_water integer check (expertise_water is null or (typeof(expertise_water) = 'integer' and expertise_water >= 0)),
    expertise_nobility integer check (expertise_nobility is null or (typeof(expertise_nobility) = 'integer' and expertise_nobility >= 0)),
    expertise_church integer check (expertise_church is null or (typeof(expertise_church) = 'integer' and expertise_church >= 0)),
    expertise_infernal integer check (expertise_infernal is null or (typeof(expertise_infernal) = 'integer' and expertise_infernal >= 0)),
    expertise_might integer check (expertise_might is null or (typeof(expertise_might) = 'integer' and expertise_might >= 0)),
    expertise_popularity integer check (expertise_popularity is null or (typeof(expertise_popularity) = 'integer' and expertise_popularity >= 0)),
    affinity_fire integer check (affinity_fire is null or (typeof(affinity_fire) = 'integer' and affinity_fire >= 0)),
    affinity_earth integer check (affinity_earth is null or (typeof(affinity_earth) = 'integer' and affinity_earth >= 0)),
    affinity_air integer check (affinity_air is null or (typeof(affinity_air) = 'integer' and affinity_air >= 0)),
    affinity_water integer check (affinity_water is null or (typeof(affinity_water) = 'integer' and affinity_water >= 0)),
    affinity_nobility integer check (affinity_nobility is null or (typeof(affinity_nobility) = 'integer' and affinity_nobility >= 0)),
    affinity_church integer check (affinity_church is null or (typeof(affinity_church) = 'integer' and affinity_church >= 0)),
    affinity_infernal integer check (affinity_infernal is null or (typeof(affinity_infernal) = 'integer' and affinity_infernal >= 0)),
    affinity_might integer check (affinity_might is null or (typeof(affinity_might) = 'integer' and affinity_might >= 0)),
    affinity_popularity integer check (affinity_popularity is null or (typeof(affinity_popularity) = 'integer' and affinity_popularity >= 0)),
    income_ephemeral_wis integer check (income_ephemeral_wis is null or (typeof(income_ephemeral_wis) = 'integer' and income_ephemeral_wis >= 0)),
    income_ephemeral_wealth integer check (income_ephemeral_wealth is null or (typeof(income_ephemeral_wealth) = 'integer' and income_ephemeral_wealth >= 0)),
    resource_wis integer check (resource_wis is null or (typeof(resource_wis) = 'integer' and resource_wis >= 0)),
    resource_wealth integer check (resource_wealth is null or (typeof(resource_wealth) = 'integer' and resource_wealth >= 0)),
    upkeep_wis integer check (upkeep_wis is null or (typeof(upkeep_wis) = 'integer' and upkeep_wis >= 0)),
    upkeep_wealth integer check (upkeep_wealth is null or (typeof(upkeep_wealth) = 'integer' and upkeep_wealth >= 0)),
    eligibility_fire integer check (eligibility_fire is null or (typeof(eligibility_fire) = 'integer' and eligibility_fire >= 0)),
    eligibility_earth integer check (eligibility_earth is null or (typeof(eligibility_earth) = 'integer' and eligibility_earth >= 0)),
    eligibility_air integer check (eligibility_air is null or (typeof(eligibility_air) = 'integer' and eligibility_air >= 0)),
    eligibility_water integer check (eligibility_water is null or (typeof(eligibility_water) = 'integer' and eligibility_water >= 0)),
    eligibility_nobility integer check (eligibility_nobility is null or (typeof(eligibility_nobility) = 'integer' and eligibility_nobility >= 0)),
    eligibility_church integer check (eligibility_church is null or (typeof(eligibility_church) = 'integer' and eligibility_church >= 0)),
    eligibility_infernal integer check (eligibility_infernal is null or (typeof(eligibility_infernal) = 'integer' and eligibility_infernal >= 0)),
    eligibility_might integer check (eligibility_might is null or (typeof(eligibility_might) = 'integer' and eligibility_might >= 0)),
    eligibility_popularity integer check (eligibility_popularity is null or (typeof(eligibility_popularity) = 'integer' and eligibility_popularity >= 0)),
    support_fire integer check (support_fire is null or (typeof(support_fire) = 'integer' and support_fire >= 0)),
    support_earth integer check (support_earth is null or (typeof(support_earth) = 'integer' and support_earth >= 0)),
    support_air integer check (support_air is null or (typeof(support_air) = 'integer' and support_air >= 0)),
    support_water integer check (support_water is null or (typeof(support_water) = 'integer' and support_water >= 0)),
    support_nobility integer check (support_nobility is null or (typeof(support_nobility) = 'integer' and support_nobility >= 0)),
    support_church integer check (support_church is null or (typeof(support_church) = 'integer' and support_church >= 0)),
    support_infernal integer check (support_infernal is null or (typeof(support_infernal) = 'integer' and support_infernal >= 0)),
    support_might integer check (support_might is null or (typeof(support_might) = 'integer' and support_might >= 0)),
    support_popularity integer check (support_popularity is null or (typeof(support_popularity) = 'integer' and support_popularity >= 0)),
    categories text not null default '[]' check (json_valid(categories) and json_type(categories) = 'array'),
    effects text not null default '[]' check (json_valid(effects) and json_type(effects) = 'array'),
    abilities text not null default '[]' check (json_valid(abilities) and json_type(abilities) = 'array'),
    challenges text not null default '[]' check (json_valid(challenges) and json_type(challenges) = 'array'),
    primary key (catalogue_version, design_id),
    unique (catalogue_version, sort_order)
);

-- Original playable-v1 catalogue, in its original deterministic order.
insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, cost_wis, affinity_fire, income_ephemeral_wis, categories, effects)
select 'playable-v1', 'chantry_fire', 'playable-v1', 0, 'Fire Chantry', 'action', 'fire', 'durable', 'discard', 0, 1, 1, 1, '["location"]', '["chantry_income"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'chantry_fire');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'element_fire', 'playable-v1', 1, 'Fire Element', 'action', 'fire', 'instant', 'discard', 0, '["element"]', '["element_support"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'element_fire');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, cost_wis, affinity_earth, income_ephemeral_wis, categories, effects)
select 'playable-v1', 'chantry_earth', 'playable-v1', 2, 'Earth Chantry', 'action', 'earth', 'durable', 'discard', 0, 1, 1, 1, '["location"]', '["chantry_income"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'chantry_earth');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'element_earth', 'playable-v1', 3, 'Earth Element', 'action', 'earth', 'instant', 'discard', 0, '["element"]', '["element_support"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'element_earth');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, cost_wis, affinity_air, income_ephemeral_wis, categories, effects)
select 'playable-v1', 'chantry_air', 'playable-v1', 4, 'Air Chantry', 'action', 'air', 'durable', 'discard', 0, 1, 1, 1, '["location"]', '["chantry_income"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'chantry_air');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'element_air', 'playable-v1', 5, 'Air Element', 'action', 'air', 'instant', 'discard', 0, '["element"]', '["element_support"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'element_air');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, cost_wis, affinity_water, income_ephemeral_wis, categories, effects)
select 'playable-v1', 'chantry_water', 'playable-v1', 6, 'Water Chantry', 'action', 'water', 'durable', 'discard', 0, 1, 1, 1, '["location"]', '["chantry_income"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'chantry_water');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, element, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'element_water', 'playable-v1', 7, 'Water Element', 'action', 'water', 'instant', 'discard', 0, '["element"]', '["element_support"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'element_water');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'research_spell', 'playable-v1', 8, 'Research Spell', 'Learn an Element, or discover a Spell from the pool selected by zero, one, or two distinct Element cards.', 'action', 'instant', 'discard', 0, 1, '["activity"]', '["research_spell"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'research_spell');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, categories)
select 'playable-v1', 'research_artifact', 'playable-v1', 9, 'Research Artifact (placeholder)', 'unavailable', 'instant', 'discard', 1, '["activity"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'research_artifact');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, resource_wis, categories, effects)
select 'playable-v1', 'basic_wis', 'playable-v1', 10, 'Basic Wis', 'action', 'instant', 'discard', 0, 1, '["resource"]', '["resource"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'basic_wis');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, resource_wealth, categories, effects)
select 'playable-v1', 'basic_wealth', 'playable-v1', 11, 'Basic Wealth', 'action', 'instant', 'discard', 0, 1, '["resource"]', '["resource"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'basic_wealth');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, text, timing, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'arcane_focus', 'playable-v1', 12, 'Arcane Focus', 'Support: 1 Wis and one chosen Affinity; serve as Challenger. Alternatively pay 2 Wis to heal 1, or support an Attack or Defence with 1 Might.', 'action_or_support', 'instant', 'discard', 0, '["activity"]', '["arcane_focus"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'arcane_focus');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, text, timing, persistence, disposal, placeholder, categories, effects, challenges)
select 'playable-v1', 'investigate', 'playable-v1', 13, 'Investigate', 'Challenge: reserve one eligible rumour using this attempt''s final domain totals.', 'action', 'instant', 'discard', 0, '["activity"]', '["investigate"]', '[{"id":"investigate"}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'investigate');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, text, timing, victory_path, persistence, disposal, placeholder, categories, effects)
select 'playable-v1', 'subjugate_rival', 'playable-v1', 14, 'Subjugate Rival', 'Initiate your one Attack this turn with ready, eligible Might Minions against your Target.', 'action', 'domination', 'instant', 'discard', 0, '["activity"]', '["subjugate_rival"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'subjugate_rival');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, victory_path, persistence, disposal, placeholder, categories)
select 'playable-v1', 'victory_arcane', 'playable-v1', 15, 'Arcane (placeholder)', 'unavailable', 'arcane', 'instant', 'discard', 1, '["placeholder"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'victory_arcane');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, victory_path, persistence, disposal, placeholder, categories)
select 'playable-v1', 'victory_influence', 'playable-v1', 16, 'Influence (placeholder)', 'unavailable', 'influence', 'instant', 'discard', 1, '["placeholder"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'victory_influence');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, victory_path, persistence, disposal, placeholder, categories)
select 'playable-v1', 'victory_council', 'playable-v1', 17, 'Council (placeholder)', 'unavailable', 'council', 'instant', 'discard', 1, '["placeholder"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'victory_council');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, victory_path, persistence, disposal, placeholder, categories)
select 'playable-v1', 'victory_fame', 'playable-v1', 18, 'Fame (placeholder)', 'unavailable', 'fame', 'instant', 'discard', 1, '["placeholder"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'victory_fame');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, expertise_fire, categories, effects)
select 'playable-v1', 'vera_coalheart', 'playable-v1', 19, 'Vera Coalheart', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'vera_coalheart');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, expertise_might, categories, effects)
select 'playable-v1', 'bastian_redhand', 'playable-v1', 20, 'Bastian Redhand', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'bastian_redhand');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, expertise_infernal, categories, effects)
select 'playable-v1', 'orla_nine_cinders', 'playable-v1', 21, 'Orla Nine-Cinders', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'orla_nine_cinders');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, expertise_earth, categories, effects)
select 'playable-v1', 'hedda_stonewake', 'playable-v1', 22, 'Hedda Stonewake', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'hedda_stonewake');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, expertise_nobility, categories, effects)
select 'playable-v1', 'leontine_grange', 'playable-v1', 23, 'Leontine Grange', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'leontine_grange');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, expertise_might, categories, effects)
select 'playable-v1', 'rurik_ironvale', 'playable-v1', 24, 'Rurik Ironvale', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rurik_ironvale');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, expertise_air, categories, effects)
select 'playable-v1', 'mina_cloudscript', 'playable-v1', 25, 'Mina Cloudscript', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'mina_cloudscript');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, expertise_might, categories, effects)
select 'playable-v1', 'garrick_stormward', 'playable-v1', 26, 'Garrick Stormward', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'garrick_stormward');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, expertise_popularity, categories, effects)
select 'playable-v1', 'elsie_many_tongues', 'playable-v1', 27, 'Elsie Many-Tongues', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'elsie_many_tongues');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, expertise_water, categories, effects)
select 'playable-v1', 'nadia_rivermark', 'playable-v1', 28, 'Nadia Rivermark', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'nadia_rivermark');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, expertise_church, categories, effects)
select 'playable-v1', 'sister_agnes_reed', 'playable-v1', 29, 'Sister Agnes Reed', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'sister_agnes_reed');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, expertise_might, categories, effects)
select 'playable-v1', 'tamsin_lockward', 'playable-v1', 30, 'Tamsin Lockward', 'action', 'durable', 'discard', 0, 1, 1, 1, '["minion"]', '["starting_minion"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'tamsin_lockward');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, resource_wis, categories, effects)
select 'playable-v1', 'ephemeral_wis', 'playable-v1', 31, 'Ephemeral wis', 'action', 'instant', 'destroy', 0, 0, 1, '["resource"]', '["ephemeral_wis"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'ephemeral_wis');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, resource_wealth, categories, effects)
select 'playable-v1', 'ephemeral_wealth', 'playable-v1', 32, 'Ephemeral wealth', 'action', 'instant', 'destroy', 0, 0, 1, '["resource"]', '["ephemeral_wealth"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'ephemeral_wealth');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, categories, effects)
select 'playable-v1', 'flame_lash', 'playable-v1', 33, 'fire', 'Flame Lash', 'During your Attack, add 2 attacking Might.', 'attack_reaction', 'instant', 'discard', 0, 1, 1, '["spell"]', '["flame_lash"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'flame_lash');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, categories, effects)
select 'playable-v1', 'firewall', 'playable-v1', 34, 'fire', 'Firewall', 'Discard one committed attacking Minion.', 'defence_reaction', 'instant', 'discard', 0, 2, 2, '["spell"]', '["firewall"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'firewall');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, support_fire, categories, effects)
select 'playable-v1', 'flamecraft', 'playable-v1', 35, 'fire', 'Flamecraft', 'Challenge support: 2 Fire.', 'support', 'instant', 'discard', 0, 0, 1, 2, '["spell"]', '["flamecraft"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'flamecraft');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, affinity_fire, categories, effects)
select 'playable-v1', 'inner_flame', 'playable-v1', 36, 'fire', 'Inner Flame', 'Provide 1 Fire Affinity.', 'action', 'durable', 'discard', 0, 2, 1, 1, '["spell"]', '["inner_flame"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'inner_flame');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, categories, effects)
select 'playable-v1', 'rockslide', 'playable-v1', 37, 'earth', 'Rockslide', 'Each committed attacker contributes 1 additional Might.', 'attack_reaction', 'instant', 'discard', 0, 1, 1, '["spell"]', '["rockslide"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rockslide');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, categories, effects)
select 'playable-v1', 'stoneward', 'playable-v1', 38, 'earth', 'Stoneward', 'Add 2 defensive Might.', 'defence_reaction', 'instant', 'discard', 0, 0, 1, '["spell"]', '["stoneward"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'stoneward');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, support_earth, categories, effects)
select 'playable-v1', 'shape_stone', 'playable-v1', 39, 'earth', 'Shape Stone', 'Challenge support: 2 Earth.', 'support', 'instant', 'discard', 0, 0, 1, 2, '["spell"]', '["shape_stone"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'shape_stone');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, upkeep_wis, categories, effects)
select 'playable-v1', 'rampart', 'playable-v1', 40, 'earth', 'Rampart', 'Add 2 defensive Might. Upkeep: 1 Wis.', 'action', 'durable', 'discard', 0, 2, 2, 1, '["spell"]', '["rampart"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rampart');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, categories, effects)
select 'playable-v1', 'dislocate', 'playable-v1', 41, 'air', 'Dislocate', 'Return one committed defender to its owner''s hand.', 'attack_reaction', 'instant', 'discard', 0, 1, 2, '["spell"]', '["dislocate"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'dislocate');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, categories, effects)
select 'playable-v1', 'deflect', 'playable-v1', 42, 'air', 'Deflect', 'Return one committed attacker to its owner''s hand.', 'defence_reaction', 'instant', 'discard', 0, 1, 1, '["spell"]', '["deflect"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'deflect');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, categories, effects)
select 'playable-v1', 'keen_insight', 'playable-v1', 43, 'air', 'Keen Insight', 'Challenge support: 1 in one chosen Affinity.', 'support', 'instant', 'discard', 0, 0, 1, '["spell"]', '["keen_insight"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'keen_insight');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, categories, effects, abilities)
select 'playable-v1', 'whispering_wind', 'playable-v1', 44, 'air', 'Whispering Wind', 'At Start.Effects, optionally draw one card, then discard one.', 'action', 'durable', 'discard', 0, 2, 2, '["spell"]', '["whispering_wind"]', '[{"id":"whispering_wind_trigger","timing":"start_effects","exhaust":false,"optional":true,"effect":"whispering_wind_trigger","text":"Draw one card, then discard one."}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'whispering_wind');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, categories, effects)
select 'playable-v1', 'deep_freeze', 'playable-v1', 45, 'water', 'Deep Freeze', 'Exhaust one ready neighbouring Minion; it skips its next Recovery.', 'action', 'instant', 'discard', 0, 2, 2, '["spell"]', '["deep_freeze"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'deep_freeze');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, categories, effects)
select 'playable-v1', 'moat', 'playable-v1', 46, 'water', 'Moat', 'Up to two committed attackers skip their next Recovery.', 'defence_reaction', 'instant', 'discard', 0, 0, 1, '["spell"]', '["moat"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'moat');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, categories, effects)
select 'playable-v1', 'flowing_form', 'playable-v1', 47, 'water', 'Flowing Form', 'Challenge support: 1 in each of two distinct chosen Affinities.', 'support', 'instant', 'discard', 0, 0, 1, '["spell"]', '["flowing_form"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'flowing_form');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_water, categories, effects, abilities)
select 'playable-v1', 'spell_healing_spring', 'playable-v1', 48, 'water', 'Healing Spring', 'Action: pay 3 Wis and exhaust this to restore 1 Integrity, maximum 3.', 'action', 'durable', 'discard', 0, 2, 2, '["spell"]', '["spell_healing_spring"]', '[{"id":"heal","timing":"action","cost":{"wis":3},"exhaust":true,"optional":false,"effect":"heal_activation","text":"Restore 1 Integrity, maximum 3."}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'spell_healing_spring');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'counterspell', 'playable-v1', 49, 'unaligned', 'Counterspell', 'Counter one pending Spell effect.', 'spell_reaction', 'instant', 'discard', 0, 1, '["spell"]', '["counterspell"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'counterspell');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'study', 'playable-v1', 50, 'unaligned', 'Study', 'Draw two, then discard one.', 'action', 'instant', 'discard', 0, 1, '["spell"]', '["study"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'study');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'unbound_will', 'playable-v1', 51, 'unaligned', 'Unbound Will', 'Permit a Challenge without a Minion or Arcane Focus; supply no capacity.', 'support', 'instant', 'discard', 0, 0, '["spell"]', '["unbound_will"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'unbound_will');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'mnemonic_seal', 'playable-v1', 52, 'unaligned', 'Mnemonic Seal', 'Increase maximum hand size by one.', 'action', 'durable', 'discard', 0, 2, '["spell"]', '["mnemonic_seal"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'mnemonic_seal');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, play_requirement_fire, categories, effects)
select 'playable-v1', 'molten_armour', 'playable-v1', 53, 'earth_fire', 'Molten Armour', 'Add 2 defensive Might; discard one committed attacker with Might Expertise at most 1.', 'defence_reaction', 'instant', 'discard', 0, 2, 1, 1, '["spell"]', '["molten_armour"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'molten_armour');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, play_requirement_fire, affinity_earth, affinity_fire, categories, effects)
select 'playable-v1', 'forgeheart', 'playable-v1', 54, 'earth_fire', 'Forgeheart', 'Provide Fire 1 and Earth 1.', 'action', 'durable', 'discard', 0, 2, 1, 1, 1, 1, '["spell"]', '["forgeheart"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'forgeheart');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_fire, categories, effects)
select 'playable-v1', 'wildfire', 'playable-v1', 55, 'air_fire', 'Wildfire', 'Discard one committed defender.', 'attack_reaction', 'instant', 'discard', 0, 2, 1, 1, '["spell"]', '["wildfire"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'wildfire');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_fire, categories, effects)
select 'playable-v1', 'flash_of_insight', 'playable-v1', 56, 'air_fire', 'Flash of Insight', 'Draw two, then discard or destroy one card from hand.', 'action', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["flash_of_insight"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'flash_of_insight');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, play_requirement_water, categories, effects)
select 'playable-v1', 'scalding_mist', 'playable-v1', 57, 'fire_water', 'Scalding Mist', 'Discard one opposing committed Minion with Might Expertise at most 1; up to one other opponent skips its next Recovery.', 'combat_reaction', 'instant', 'discard', 0, 2, 1, 1, '["spell"]', '["scalding_mist"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'scalding_mist');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_fire, play_requirement_water, categories, effects)
select 'playable-v1', 'purifying_steam', 'playable-v1', 58, 'fire_water', 'Purifying Steam', 'Additional cost: discard one hand card. Restore 1 Integrity, maximum 3.', 'action', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["purifying_steam"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'purifying_steam');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_earth, categories, effects)
select 'playable-v1', 'sandstorm', 'playable-v1', 59, 'air_earth', 'Sandstorm', 'Opposing committed Minions each contribute 1 less Might, minimum zero.', 'combat_reaction', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["sandstorm"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'sandstorm');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_earth, categories, effects)
select 'playable-v1', 'wind_road', 'playable-v1', 60, 'air_earth', 'Wind Road', 'Reduce Minion Wis play costs by 1, minimum zero.', 'action', 'durable', 'discard', 0, 2, 1, 1, '["spell"]', '["wind_road"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'wind_road');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, play_requirement_water, categories, effects)
select 'playable-v1', 'quagmire', 'playable-v1', 61, 'earth_water', 'Quagmire', 'Up to two attackers each contribute 1 less Might and skip their next Recovery.', 'defence_reaction', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["quagmire"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'quagmire');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_earth, play_requirement_water, income_ephemeral_wealth, categories, effects)
select 'playable-v1', 'verdant_spring', 'playable-v1', 62, 'earth_water', 'Verdant Spring', 'Generate one ephemeral Wealth at Start.Income.', 'action', 'durable', 'discard', 0, 2, 1, 1, 1, '["spell"]', '["verdant_spring"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'verdant_spring');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_water, categories, effects)
select 'playable-v1', 'freezing_fog', 'playable-v1', 63, 'air_water', 'Freezing Fog', 'Return one attacker with Might Expertise at most 1 to hand; up to one other attacker skips its next Recovery.', 'defence_reaction', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["freezing_fog"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'freezing_fog');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, pool, name, text, timing, persistence, disposal, placeholder, cost_wis, play_requirement_air, play_requirement_water, categories, effects)
select 'playable-v1', 'scrying_pool', 'playable-v1', 64, 'air_water', 'Scrying Pool', 'Privately inspect up to the top three draw cards, take one, and reorder the rest; do not reshuffle.', 'action', 'instant', 'discard', 0, 1, 1, 1, '["spell"]', '["scrying_pool"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'scrying_pool');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, text, timing, persistence, disposal, placeholder, cost_wis, categories, effects)
select 'playable-v1', 'pupil', 'playable-v1', 65, 'Pupil', 'On recruitment choose Expertise 1 uniformly from the eight domains other than the drafted Chantry''s Element. This trait is permanent.', 'action', 'durable', 'discard', 0, 1, '["minion"]', '["pupil"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'pupil');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, expertise_infernal, categories, effects)
select 'playable-v1', 'witch', 'playable-v1', 66, 'Witch', 'action', 'durable', 'discard', 0, 1, 1, '["minion"]', '["witch"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'witch');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, income_ephemeral_wis, categories, effects)
select 'playable-v1', 'leyline', 'playable-v1', 67, 'Leyline', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["leyline"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'leyline');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_nobility, categories, effects)
select 'playable-v1', 'court_office', 'playable-v1', 68, 'Court Office', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["court_office"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'court_office');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_popularity, categories, effects)
select 'playable-v1', 'village_hall', 'playable-v1', 69, 'Village Hall', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["village_hall"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'village_hall');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_earth, categories, effects)
select 'playable-v1', 'stone_circle', 'playable-v1', 70, 'Stone Circle', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["stone_circle"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'stone_circle');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_water, categories, effects)
select 'playable-v1', 'healing_spring', 'playable-v1', 71, 'Healing Spring', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["healing_spring"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'healing_spring');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, expertise_church, affinity_church, categories, effects)
select 'playable-v1', 'abbey_liaison', 'playable-v1', 72, 'Abbey Liaison', 'action', 'durable', 'discard', 0, 2, 1, 1, '["minion"]', '["abbey_liaison"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'abbey_liaison');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_infernal, categories, effects)
select 'playable-v1', 'crossroads_shrine', 'playable-v1', 73, 'Crossroads Shrine', 'action', 'durable', 'discard', 0, 1, 1, '["location"]', '["crossroads_shrine"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'crossroads_shrine');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, expertise_might, affinity_might, categories, effects)
select 'playable-v1', 'caravan_guard', 'playable-v1', 74, 'Caravan Guard', 'action', 'durable', 'discard', 0, 2, 1, 1, '["minion"]', '["caravan_guard"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'caravan_guard');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_fire, affinity_infernal, categories, effects)
select 'playable-v1', 'ember_shrine', 'playable-v1', 75, 'Ember Shrine', 'action', 'durable', 'discard', 0, 2, 1, 1, '["location"]', '["ember_shrine"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'ember_shrine');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, affinity_earth, affinity_nobility, income_ephemeral_wealth, categories, effects)
select 'playable-v1', 'chartered_quarry', 'playable-v1', 76, 'Chartered Quarry', 'action', 'durable', 'discard', 0, 3, 1, 1, 1, '["location"]', '["chartered_quarry"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'chartered_quarry');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, persistence, disposal, placeholder, cost_wis, expertise_church, expertise_water, affinity_church, categories, effects)
select 'playable-v1', 'reedbound_healer', 'playable-v1', 77, 'Reedbound Healer', 'action', 'durable', 'discard', 0, 2, 1, 1, 1, '["minion"]', '["reedbound_healer"]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'reedbound_healer');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_promising_pupil', 'playable-v1', 78, 'Promising Pupil', 'start', 2, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_promising_pupil"]', '[{"id":"accept_pupil","cost":{"wis":3},"requiresfocus":true,"rewards":["pupil"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_promising_pupil');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, abilities)
select 'playable-v1', 'rumour_solar_eclipse', 'playable-v1', 79, 'Solar Eclipse', 'start', 1, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_solar_eclipse"]', '[{"id":"solar_eclipse","timing":"start_effects","exhaust":false,"optional":false,"effect":"solar_eclipse","text":"Gain 2 Wis, Fire 1 and Infernal 1 until this turn ends."}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_solar_eclipse');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, abilities, challenges)
select 'playable-v1', 'rumour_burn_the_witch', 'playable-v1', 80, 'Burn the Witch', 'start', 3, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_burn_the_witch"]', '[{"id":"burn_failure","timing":"expiry_failure","exhaust":false,"optional":false,"effect":"burn_failure","text":"If no Challenge succeeded before expiry, discard one chosen Asset; defer to next Start.Effects when this turn''s End Phase is skipped."}]', '[{"id":"arrest_witch","requirements":{"might":1}},{"id":"recruit_witch","cost":{"wis":2},"rewards":["witch"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_burn_the_witch');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_new_leyline', 'playable-v1', 81, 'New Leyline', 'start', 2, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_new_leyline"]', '[{"id":"claim_leyline","cost":{"wis":2},"rewards":["leyline"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_new_leyline');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_full_moon', 'playable-v1', 82, 'Full Moon', 'start', 1, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_full_moon"]', '[{"rewardduration":2,"id":"moonlit_hunt","requirements":{"might":1},"rewards":["basic_wealth"],"affinityreward":{"infernal":1,"water":1}}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_full_moon');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_wandering_merchant', 'playable-v1', 83, 'Wandering Merchant', 'start', 2, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_wandering_merchant"]', '[{"id":"strike_bargain","cost":{"wealth":1},"rewards":["basic_wis","basic_wealth"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_wandering_merchant');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_lost_courier', 'playable-v1', 84, 'Lost Courier', 'start', 2, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_lost_courier"]', '[{"id":"follow_messages","requirements":{"air":1},"rewards":["court_office"]},{"rewardduration":2,"id":"escort_courier","requirements":{"might":1},"rewards":["basic_wealth"],"affinityreward":{"nobility":1}}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_lost_courier');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, categories, effects, challenges)
select 'playable-v1', 'rumour_village_petition', 'playable-v1', 85, 'Village Petition', 'start', 2, 'duration', 'mill', 0, 0, '["rumour"]', '["rumour_village_petition"]', '[{"id":"rally_villagers","requirements":{"popularity":1},"cost":{"wis":1},"rewards":["village_hall"]},{"id":"repair_well","requirements":{"earth":1},"cost":{"wis":1},"rewards":["village_hall"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_village_petition');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_fire, categories, effects, challenges)
select 'playable-v1', 'rumour_restless_kiln', 'playable-v1', 86, 'Restless Kiln', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_restless_kiln"]', '[{"id":"tame_furnace","requirements":{"fire":2},"rewards":["basic_wealth","basic_wealth"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_restless_kiln');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_earth, categories, effects, challenges)
select 'playable-v1', 'rumour_boundary_stones', 'playable-v1', 87, 'Buried Boundary Stones', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_boundary_stones"]', '[{"id":"raise_stones","requirements":{"earth":2},"rewards":["stone_circle"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_boundary_stones');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_air, categories, effects, challenges)
select 'playable-v1', 'rumour_whispers', 'playable-v1', 88, 'Whispers on the Wind', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_whispers"]', '[{"rewardduration":2,"id":"gather_whispers","requirements":{"air":2},"rewards":["basic_wis"],"affinityreward":{"popularity":1}}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_whispers');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_water, categories, effects, challenges)
select 'playable-v1', 'rumour_healing_spring', 'playable-v1', 89, 'The Healing Spring', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_healing_spring"]', '[{"id":"restore_spring","requirements":{"water":2},"rewards":["healing_spring"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_healing_spring');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_nobility, categories, effects, challenges)
select 'playable-v1', 'rumour_inheritance', 'playable-v1', 90, 'Disputed Inheritance', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_inheritance"]', '[{"id":"witness_settlement","requirements":{"nobility":2},"cost":{"wis":1},"rewards":["court_office"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_inheritance');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_church, categories, effects, challenges)
select 'playable-v1', 'rumour_infirmary', 'playable-v1', 91, 'Empty Abbey Infirmary', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_infirmary"]', '[{"id":"organise_relief","requirements":{"church":1},"cost":{"wis":1},"rewards":["abbey_liaison"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_infirmary');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_infernal, categories, effects, challenges)
select 'playable-v1', 'rumour_crossroads', 'playable-v1', 92, 'Voices at the Crossroads', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_crossroads"]', '[{"id":"make_offering","requirements":{"infernal":1},"cost":{"wealth":1},"rewards":["crossroads_shrine"]},{"id":"learn_names","requirements":{"infernal":2},"rewards":["crossroads_shrine"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_crossroads');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_might, categories, effects, challenges)
select 'playable-v1', 'rumour_caravan', 'playable-v1', 93, 'Caravan under Escort', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_caravan"]', '[{"id":"guard_road","requirements":{"might":2},"rewards":["basic_wealth","basic_wealth"]},{"id":"hire_guard","requirements":{"might":1},"cost":{"wealth":1},"rewards":["caravan_guard"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_caravan');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_popularity, categories, effects, challenges)
select 'playable-v1', 'rumour_village_fair', 'playable-v1', 94, 'Village Fair', 'start', 2, 'duration', 'mill', 0, 0, 1, '["rumour"]', '["rumour_village_fair"]', '[{"id":"win_crowd","requirements":{"popularity":2},"rewards":["village_hall","basic_wealth"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_village_fair');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_fire, eligibility_infernal, categories, effects, challenges)
select 'playable-v1', 'rumour_gallows', 'playable-v1', 95, 'Embers beneath the Gallows', 'start', 3, 'duration', 'mill', 0, 0, 1, 1, '["rumour"]', '["rumour_gallows"]', '[{"id":"bind_embers","requirements":{"fire":2,"infernal":2},"cost":{"wis":1},"rewards":["ember_shrine"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_gallows');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_earth, eligibility_nobility, categories, effects, challenges)
select 'playable-v1', 'rumour_quarry', 'playable-v1', 96, 'The Baron''s Forgotten Quarry', 'start', 3, 'duration', 'mill', 0, 0, 1, 1, '["rumour"]', '["rumour_quarry"]', '[{"id":"reopen_quarry","requirements":{"earth":2,"nobility":1},"cost":{"wealth":1},"rewards":["chartered_quarry"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_quarry');

insert into card_designs (catalogue_version, design_id, rules_version, sort_order, name, timing, duration, persistence, disposal, placeholder, cost_wis, eligibility_church, eligibility_water, categories, effects, challenges)
select 'playable-v1', 'rumour_crypt', 'playable-v1', 97, 'The Abbey''s Flooded Crypt', 'start', 3, 'duration', 'mill', 0, 0, 1, 1, '["rumour"]', '["rumour_crypt"]', '[{"id":"recover_records","requirements":{"church":1,"water":2},"cost":{"wis":1},"rewards":["reedbound_healer"]}]'
where not exists (select 1 from card_designs where catalogue_version = 'playable-v1' and design_id = 'rumour_crypt');

-- A version becomes immutable when its first match is saved, including historical matches.
create trigger if not exists card_designs_frozen_insert
before insert on card_designs
when exists (select 1 from matches where case when json_valid(snapshot) then json_extract(snapshot, '$.catalogueversion') end = new.catalogue_version)
begin
    select raise(abort, 'catalogue version is used by a match and cannot be changed');
end;

create trigger if not exists card_designs_frozen_update
before update on card_designs
when exists (select 1 from matches where case when json_valid(snapshot) then json_extract(snapshot, '$.catalogueversion') end = old.catalogue_version) or exists (select 1 from matches where case when json_valid(snapshot) then json_extract(snapshot, '$.catalogueversion') end = new.catalogue_version)
begin
    select raise(abort, 'catalogue version is used by a match and cannot be changed');
end;

create trigger if not exists card_designs_frozen_delete
before delete on card_designs
when exists (select 1 from matches where case when json_valid(snapshot) then json_extract(snapshot, '$.catalogueversion') end = old.catalogue_version)
begin
    select raise(abort, 'catalogue version is used by a match and cannot be changed');
end;

commit;
