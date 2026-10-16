-- Insight (enchant 8216): the 70170 hotfixes moved its buff 1299796 from A_MOD_PERCENT_STAT Spirit
-- (aura 80, misc 4) to A_MOD_TOTAL_STAT_PERCENTAGE with neither misc value set, which names no stat
-- (or Strength, read as a stat index). Its tooltip still says "Increases your Spirit by $s1%", so the
-- 70124 row is put back. The WHERE clause stops matching once the client states a stat itself, which
-- client 70425 does (stat mask 16, Spirit).
UPDATE SpellEffect SET EffectAura = 80, EffectMiscValue = json_set(EffectMiscValue, '$[0]', 4)
WHERE SpellID = 1299796 AND EffectAura = 137 AND EffectMiscValue_0 = 0 AND EffectMiscValue_1 = 0;
