// Code generated by "stringer -output=rule_name_string.go -type=RuleName rule_name.go rule_name.og.go"; DO NOT EDIT.

package opt

import "strconv"

const _RuleName_name = "InvalidRuleNameSimplifyRootOrderingPruneRootColsNumManualRuleNamesEliminateEmptyAndEliminateEmptyOrEliminateSingletonAndOrSimplifyAndSimplifyOrSimplifyFiltersFoldNullAndOrNegateComparisonEliminateNotNegateAndNegateOrExtractRedundantClauseExtractRedundantSubclauseCommuteVarInequalityCommuteConstInequalityNormalizeCmpPlusConstNormalizeCmpMinusConstNormalizeCmpConstMinusNormalizeTupleEqualityFoldNullComparisonLeftFoldNullComparisonRightFoldIsNullFoldNonNullIsNullFoldIsNotNullFoldNonNullIsNotNullCommuteNullIsDecorrelateJoinTryDecorrelateSelectTryDecorrelateProjectTryDecorrelateProjectSelectTryDecorrelateScalarGroupByHoistSelectExistsHoistSelectNotExistsHoistSelectSubqueryHoistProjectSubqueryHoistJoinSubqueryHoistValuesSubqueryNormalizeAnyFilterNormalizeNotAnyFilterEliminateDistinctEliminateGroupByProjectPushSelectIntoInlinableProjectInlineProjectInProjectEnsureJoinFiltersAndEnsureJoinFiltersPushFilterIntoJoinLeftPushFilterIntoJoinRightSimplifyLeftJoinSimplifyRightJoinEliminateSemiJoinEliminateAntiJoinEliminateJoinNoColsLeftEliminateJoinNoColsRightEliminateLimitPushLimitIntoProjectPushOffsetIntoProjectEliminateMax1RowFoldPlusZeroFoldZeroPlusFoldMinusZeroFoldMultOneFoldOneMultFoldDivOneInvertMinusEliminateUnaryMinusFoldUnaryMinusSimplifyLimitOrderingSimplifyOffsetOrderingSimplifyGroupByOrderingSimplifyRowNumberOrderingSimplifyExplainOrderingEliminateProjectEliminateProjectProjectPruneProjectColsPruneScanColsPruneSelectColsPruneLimitColsPruneOffsetColsPruneJoinLeftColsPruneJoinRightColsPruneAggColsPruneGroupByColsPruneValuesColsPruneRowNumberColsCommuteVarCommuteConstEliminateCoalesceSimplifyCoalesceEliminateCastFoldNullCastFoldNullUnaryFoldNullBinaryLeftFoldNullBinaryRightFoldNullInNonEmptyFoldNullInEmptyFoldNullNotInEmptyNormalizeInConstFoldInNullEliminateExistsProjectEliminateExistsGroupByEliminateSelectEnsureSelectFiltersAndEnsureSelectFiltersMergeSelectsPushSelectIntoProjectPushSelectIntoJoinLeftPushSelectIntoJoinRightMergeSelectInnerJoinPushSelectIntoGroupByRemoveNotNullConditionGenerateMergeJoinsPushLimitIntoScanPushLimitIntoLookupJoinGenerateIndexScansConstrainScanPushFilterIntoLookupJoinNoRemainderPushFilterIntoLookupJoinConstrainLookupJoinIndexScanNumRuleNames"

var _RuleName_index = [...]uint16{0, 15, 35, 48, 66, 83, 99, 122, 133, 143, 158, 171, 187, 199, 208, 216, 238, 263, 283, 305, 326, 348, 370, 392, 414, 437, 447, 464, 477, 497, 510, 525, 545, 566, 593, 620, 637, 657, 676, 696, 713, 732, 750, 771, 788, 811, 841, 863, 883, 900, 922, 945, 961, 978, 995, 1012, 1035, 1059, 1073, 1093, 1114, 1130, 1142, 1154, 1167, 1178, 1189, 1199, 1210, 1229, 1243, 1264, 1286, 1309, 1334, 1357, 1373, 1396, 1412, 1425, 1440, 1454, 1469, 1486, 1504, 1516, 1532, 1547, 1565, 1575, 1587, 1604, 1620, 1633, 1645, 1658, 1676, 1695, 1713, 1728, 1746, 1762, 1772, 1794, 1816, 1831, 1853, 1872, 1884, 1905, 1927, 1950, 1970, 1991, 2013, 2031, 2048, 2071, 2089, 2102, 2137, 2161, 2189, 2201}

func (i RuleName) String() string {
	if i >= RuleName(len(_RuleName_index)-1) {
		return "RuleName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleName_name[_RuleName_index[i]:_RuleName_index[i+1]]
}
