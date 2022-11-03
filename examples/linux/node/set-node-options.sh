#!/bin/sh

set -ex

NODE_OPTIONS=()

if [ -n "$ENABLE_FIPS" ] && [ $ENABLE_FIPS == "y" ]; then
  NODE_OPTIONS+=("--enable-fips")
fi

if [ -n "$NO_DEPRECATION" ] && [ $NO_DEPRECATION == "y" ]; then
  NODE_OPTIONS+=("--no-deprecation")
fi

if [ -n "$NO_WARNINGS" ] && [ $NO_WARNINGS == "y" ]; then
  NODE_OPTIONS+=("--no-warnings")
fi

if [ -n "$TRACE_SYNC_IO" ] && [ $TRACE_SYNC_IO == "y" ]; then
  NODE_OPTIONS+=("--trace-sync-io")
fi

if [ -n "$TRACE_DEPRECATION" ] && [ $TRACE_DEPRECATION == "y" ]; then
  NODE_OPTIONS+=("--trace-deprecation")
fi

if [ -n "$TRACE_WARNINGS" ] && [ $TRACE_WARNINGS == "y" ]; then
  NODE_OPTIONS+=("--trace-warnings")
fi

if [ -n "$TRACK_HEAP_OBJECTS" ] && [ $TRACK_HEAP_OBJECTS == "y" ]; then
  NODE_OPTIONS+=("--track-heap-objects")
fi

if [ -n "$ZERO_FILL_BUFFERS" ] && [ $ZERO_FILL_BUFFERS == "y" ]; then
  NODE_OPTIONS+=("--zero-fill-buffers")
fi

if [ -n "$EXPERIMENTAL_EXTRAS" ] && [ $EXPERIMENTAL_EXTRAS == "y" ]; then
  NODE_OPTIONS+=("--experimental_extras")
else
  NODE_OPTIONS+=("--no-experimental_extras")
fi

if [ -n "$USE_STRICT" ] && [ $USE_STRICT == "y" ]; then
  NODE_OPTIONS+=("--use_strict")
else
  NODE_OPTIONS+=("--no-use_strict")
fi

if [ -n "$ES_STAGING" ] && [ $ES_STAGING == "y" ]; then
  NODE_OPTIONS+=("--es_staging")
else
  NODE_OPTIONS+=("--no-es_staging")
fi

if [ -n "$HARMONY" ] && [ $HARMONY == "y" ]; then
  NODE_OPTIONS+=("--harmony")
else
  NODE_OPTIONS+=("--no-harmony")
fi

if [ -n "$HARMONY_SHIPPING" ] && [ $HARMONY_SHIPPING == "y" ]; then
  NODE_OPTIONS+=("--harmony_shipping")
else
  NODE_OPTIONS+=("--no-harmony_shipping")
fi

if [ -n "$HARMONY_ARRAY_PROTOTYPE_VALUES" ] && [ $HARMONY_ARRAY_PROTOTYPE_VALUES == "y" ]; then
  NODE_OPTIONS+=("--harmony_array_prototype_values")
else
  NODE_OPTIONS+=("--no-harmony_array_prototype_values")
fi

if [ -n "$HARMONY_FUNCTION_SENT" ] && [ $HARMONY_FUNCTION_SENT == "y" ]; then
  NODE_OPTIONS+=("--harmony_function_sent")
else
  NODE_OPTIONS+=("--no-harmony_function_sent")
fi

if [ -n "$HARMONY_SHAREDARRAYBUFFER" ] && [ $HARMONY_SHAREDARRAYBUFFER == "y" ]; then
  NODE_OPTIONS+=("--harmony_sharedarraybuffer")
else
  NODE_OPTIONS+=("--no-harmony_sharedarraybuffer")
fi

if [ -n "$HARMONY_DO_EXPRESSIONS" ] && [ $HARMONY_DO_EXPRESSIONS == "y" ]; then
  NODE_OPTIONS+=("--harmony_do_expressions")
else
  NODE_OPTIONS+=("--no-harmony_do_expressions")
fi

if [ -n "$HARMONY_REGEXP_NAMED_CAPTURES" ] && [ $HARMONY_REGEXP_NAMED_CAPTURES == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_named_captures")
else
  NODE_OPTIONS+=("--no-harmony_regexp_named_captures")
fi

if [ -n "$HARMONY_REGEXP_PROPERTY" ] && [ $HARMONY_REGEXP_PROPERTY == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_property")
else
  NODE_OPTIONS+=("--no-harmony_regexp_property")
fi

if [ -n "$HARMONY_FUNCTION_TOSTRING" ] && [ $HARMONY_FUNCTION_TOSTRING == "y" ]; then
  NODE_OPTIONS+=("--harmony_function_tostring")
else
  NODE_OPTIONS+=("--no-harmony_function_tostring")
fi

if [ -n "$HARMONY_CLASS_FIELDS" ] && [ $HARMONY_CLASS_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--harmony_class_fields")
else
  NODE_OPTIONS+=("--no-harmony_class_fields")
fi

if [ -n "$HARMONY_ASYNC_ITERATION" ] && [ $HARMONY_ASYNC_ITERATION == "y" ]; then
  NODE_OPTIONS+=("--harmony_async_iteration")
else
  NODE_OPTIONS+=("--no-harmony_async_iteration")
fi

if [ -n "$HARMONY_DYNAMIC_IMPORT" ] && [ $HARMONY_DYNAMIC_IMPORT == "y" ]; then
  NODE_OPTIONS+=("--harmony_dynamic_import")
else
  NODE_OPTIONS+=("--no-harmony_dynamic_import")
fi

if [ -n "$HARMONY_PROMISE_FINALLY" ] && [ $HARMONY_PROMISE_FINALLY == "y" ]; then
  NODE_OPTIONS+=("--harmony_promise_finally")
else
  NODE_OPTIONS+=("--no-harmony_promise_finally")
fi

if [ -n "$HARMONY_REGEXP_LOOKBEHIND" ] && [ $HARMONY_REGEXP_LOOKBEHIND == "y" ]; then
  NODE_OPTIONS+=("--harmony_regexp_lookbehind")
else
  NODE_OPTIONS+=("--no-harmony_regexp_lookbehind")
fi

if [ -n "$HARMONY_RESTRICTIVE_GENERATORS" ] && [ $HARMONY_RESTRICTIVE_GENERATORS == "y" ]; then
  NODE_OPTIONS+=("--harmony_restrictive_generators")
else
  NODE_OPTIONS+=("--no-harmony_restrictive_generators")
fi

if [ -n "$HARMONY_OBJECT_REST_SPREAD" ] && [ $HARMONY_OBJECT_REST_SPREAD == "y" ]; then
  NODE_OPTIONS+=("--harmony_object_rest_spread")
else
  NODE_OPTIONS+=("--no-harmony_object_rest_spread")
fi

if [ -n "$HARMONY_TEMPLATE_ESCAPES" ] && [ $HARMONY_TEMPLATE_ESCAPES == "y" ]; then
  NODE_OPTIONS+=("--harmony_template_escapes")
else
  NODE_OPTIONS+=("--no-harmony_template_escapes")
fi

if [ -n "$FUTURE" ] && [ $FUTURE == "y" ]; then
  NODE_OPTIONS+=("--future")
else
  NODE_OPTIONS+=("--no-future")
fi

if [ -n "$ALLOCATION_SITE_PRETENURING" ] && [ $ALLOCATION_SITE_PRETENURING == "y" ]; then
  NODE_OPTIONS+=("--allocation_site_pretenuring")
else
  NODE_OPTIONS+=("--no-allocation_site_pretenuring")
fi

if [ -n "$PAGE_PROMOTION" ] && [ $PAGE_PROMOTION == "y" ]; then
  NODE_OPTIONS+=("--page_promotion")
else
  NODE_OPTIONS+=("--no-page_promotion")
fi

if [ -n "$TRACE_PRETENURING" ] && [ $TRACE_PRETENURING == "y" ]; then
  NODE_OPTIONS+=("--trace_pretenuring")
else
  NODE_OPTIONS+=("--no-trace_pretenuring")
fi

if [ -n "$TRACE_PRETENURING_STATISTICS" ] && [ $TRACE_PRETENURING_STATISTICS == "y" ]; then
  NODE_OPTIONS+=("--trace_pretenuring_statistics")
else
  NODE_OPTIONS+=("--no-trace_pretenuring_statistics")
fi

if [ -n "$TRACK_FIELDS" ] && [ $TRACK_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_fields")
else
  NODE_OPTIONS+=("--no-track_fields")
fi

if [ -n "$TRACK_DOUBLE_FIELDS" ] && [ $TRACK_DOUBLE_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_double_fields")
else
  NODE_OPTIONS+=("--no-track_double_fields")
fi

if [ -n "$TRACK_HEAP_OBJECT_FIELDS" ] && [ $TRACK_HEAP_OBJECT_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_heap_object_fields")
else
  NODE_OPTIONS+=("--no-track_heap_object_fields")
fi

if [ -n "$TRACK_COMPUTED_FIELDS" ] && [ $TRACK_COMPUTED_FIELDS == "y" ]; then
  NODE_OPTIONS+=("--track_computed_fields")
else
  NODE_OPTIONS+=("--no-track_computed_fields")
fi

if [ -n "$TRACK_FIELD_TYPES" ] && [ $TRACK_FIELD_TYPES == "y" ]; then
  NODE_OPTIONS+=("--track_field_types")
else
  NODE_OPTIONS+=("--no-track_field_types")
fi

if [ -n "$TYPE_PROFILE" ] && [ $TYPE_PROFILE == "y" ]; then
  NODE_OPTIONS+=("--type_profile")
else
  NODE_OPTIONS+=("--no-type_profile")
fi

if [ -n "$OPTIMIZE_FOR_SIZE" ] && [ $OPTIMIZE_FOR_SIZE == "y" ]; then
  NODE_OPTIONS+=("--optimize_for_size")
else
  NODE_OPTIONS+=("--no-optimize_for_size")
fi

if [ -n "$UNBOX_DOUBLE_ARRAYS" ] && [ $UNBOX_DOUBLE_ARRAYS == "y" ]; then
  NODE_OPTIONS+=("--unbox_double_arrays")
else
  NODE_OPTIONS+=("--no-unbox_double_arrays")
fi

if [ -n "$STRING_SLICES" ] && [ $STRING_SLICES == "y" ]; then
  NODE_OPTIONS+=("--string_slices")
else
  NODE_OPTIONS+=("--no-string_slices")
fi

if [ -n "$IGNITION_REO" ] && [ $IGNITION_REO == "y" ]; then
  NODE_OPTIONS+=("--ignition_reo")
else
  NODE_OPTIONS+=("--no-ignition_reo")
fi

if [ -n "$IGNITION_FILTER_EXPRESSION_POSITIONS" ] && [ $IGNITION_FILTER_EXPRESSION_POSITIONS == "y" ]; then
  NODE_OPTIONS+=("--ignition_filter_expression_positions")
else
  NODE_OPTIONS+=("--no-ignition_filter_expression_positions")
fi

if [ -n "$PRINT_BYTECODE" ] && [ $PRINT_BYTECODE == "y" ]; then
  NODE_OPTIONS+=("--print_bytecode")
else
  NODE_OPTIONS+=("--no-print_bytecode")
fi

if [ -n "$TRACE_IGNITION_CODEGEN" ] && [ $TRACE_IGNITION_CODEGEN == "y" ]; then
  NODE_OPTIONS+=("--trace_ignition_codegen")
else
  NODE_OPTIONS+=("--no-trace_ignition_codegen")
fi

if [ -n "$TRACE_IGNITION_DISPATCHES" ] && [ $TRACE_IGNITION_DISPATCHES == "y" ]; then
  NODE_OPTIONS+=("--trace_ignition_dispatches")
else
  NODE_OPTIONS+=("--no-trace_ignition_dispatches")
fi

if [ -n "$FAST_MATH" ] && [ $FAST_MATH == "y" ]; then
  NODE_OPTIONS+=("--fast_math")
else
  NODE_OPTIONS+=("--no-fast_math")
fi

if [ -n "$TRACE_ENVIRONMENT_LIVENESS" ] && [ $TRACE_ENVIRONMENT_LIVENESS == "y" ]; then
  NODE_OPTIONS+=("--trace_environment_liveness")
else
  NODE_OPTIONS+=("--no-trace_environment_liveness")
fi

if [ -n "$TRACE_STORE_ELIMINATION" ] && [ $TRACE_STORE_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--trace_store_elimination")
else
  NODE_OPTIONS+=("--no-trace_store_elimination")
fi

if [ -n "$TRACE_ALLOC" ] && [ $TRACE_ALLOC == "y" ]; then
  NODE_OPTIONS+=("--trace_alloc")
else
  NODE_OPTIONS+=("--no-trace_alloc")
fi

if [ -n "$TRACE_ALL_USES" ] && [ $TRACE_ALL_USES == "y" ]; then
  NODE_OPTIONS+=("--trace_all_uses")
else
  NODE_OPTIONS+=("--no-trace_all_uses")
fi

if [ -n "$TRACE_REPRESENTATION" ] && [ $TRACE_REPRESENTATION == "y" ]; then
  NODE_OPTIONS+=("--trace_representation")
else
  NODE_OPTIONS+=("--no-trace_representation")
fi

if [ -n "$TRACE_TRACK_ALLOCATION_SITES" ] && [ $TRACE_TRACK_ALLOCATION_SITES == "y" ]; then
  NODE_OPTIONS+=("--trace_track_allocation_sites")
else
  NODE_OPTIONS+=("--no-trace_track_allocation_sites")
fi

if [ -n "$TRACE_MIGRATION" ] && [ $TRACE_MIGRATION == "y" ]; then
  NODE_OPTIONS+=("--trace_migration")
else
  NODE_OPTIONS+=("--no-trace_migration")
fi

if [ -n "$TRACE_GENERALIZATION" ] && [ $TRACE_GENERALIZATION == "y" ]; then
  NODE_OPTIONS+=("--trace_generalization")
else
  NODE_OPTIONS+=("--no-trace_generalization")
fi

if [ -n "$PRINT_DEOPT_STRESS" ] && [ $PRINT_DEOPT_STRESS == "y" ]; then
  NODE_OPTIONS+=("--print_deopt_stress")
else
  NODE_OPTIONS+=("--no-print_deopt_stress")
fi

if [ -n "$POLYMORPHIC_INLINING" ] && [ $POLYMORPHIC_INLINING == "y" ]; then
  NODE_OPTIONS+=("--polymorphic_inlining")
else
  NODE_OPTIONS+=("--no-polymorphic_inlining")
fi

if [ -n "$USE_OSR" ] && [ $USE_OSR == "y" ]; then
  NODE_OPTIONS+=("--use_osr")
else
  NODE_OPTIONS+=("--no-use_osr")
fi

if [ -n "$ANALYZE_ENVIRONMENT_LIVENESS" ] && [ $ANALYZE_ENVIRONMENT_LIVENESS == "y" ]; then
  NODE_OPTIONS+=("--analyze_environment_liveness")
else
  NODE_OPTIONS+=("--no-analyze_environment_liveness")
fi

if [ -n "$TRACE_OSR" ] && [ $TRACE_OSR == "y" ]; then
  NODE_OPTIONS+=("--trace_osr")
else
  NODE_OPTIONS+=("--no-trace_osr")
fi

if [ -n "$INLINE_ACCESSORS" ] && [ $INLINE_ACCESSORS == "y" ]; then
  NODE_OPTIONS+=("--inline_accessors")
else
  NODE_OPTIONS+=("--no-inline_accessors")
fi

if [ -n "$INLINE_INTO_TRY" ] && [ $INLINE_INTO_TRY == "y" ]; then
  NODE_OPTIONS+=("--inline_into_try")
else
  NODE_OPTIONS+=("--no-inline_into_try")
fi

if [ -n "$CONCURRENT_RECOMPILATION" ] && [ $CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-concurrent_recompilation")
fi

if [ -n "$TRACE_CONCURRENT_RECOMPILATION" ] && [ $TRACE_CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--trace_concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-trace_concurrent_recompilation")
fi

if [ -n "$BLOCK_CONCURRENT_RECOMPILATION" ] && [ $BLOCK_CONCURRENT_RECOMPILATION == "y" ]; then
  NODE_OPTIONS+=("--block_concurrent_recompilation")
else
  NODE_OPTIONS+=("--no-block_concurrent_recompilation")
fi

if [ -n "$TURBO_SP_FRAME_ACCESS" ] && [ $TURBO_SP_FRAME_ACCESS == "y" ]; then
  NODE_OPTIONS+=("--turbo_sp_frame_access")
else
  NODE_OPTIONS+=("--no-turbo_sp_frame_access")
fi

if [ -n "$TURBO_PREPROCESS_RANGES" ] && [ $TURBO_PREPROCESS_RANGES == "y" ]; then
  NODE_OPTIONS+=("--turbo_preprocess_ranges")
else
  NODE_OPTIONS+=("--no-turbo_preprocess_ranges")
fi

if [ -n "$TRACE_TURBO" ] && [ $TRACE_TURBO == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo")
else
  NODE_OPTIONS+=("--no-trace_turbo")
fi

if [ -n "$TRACE_TURBO_GRAPH" ] && [ $TRACE_TURBO_GRAPH == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_graph")
else
  NODE_OPTIONS+=("--no-trace_turbo_graph")
fi

if [ -n "$TRACE_TURBO_TYPES" ] && [ $TRACE_TURBO_TYPES == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_types")
else
  NODE_OPTIONS+=("--no-trace_turbo_types")
fi

if [ -n "$TRACE_TURBO_SCHEDULER" ] && [ $TRACE_TURBO_SCHEDULER == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_scheduler")
else
  NODE_OPTIONS+=("--no-trace_turbo_scheduler")
fi

if [ -n "$TRACE_TURBO_REDUCTION" ] && [ $TRACE_TURBO_REDUCTION == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_reduction")
else
  NODE_OPTIONS+=("--no-trace_turbo_reduction")
fi

if [ -n "$TRACE_TURBO_TRIMMING" ] && [ $TRACE_TURBO_TRIMMING == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_trimming")
else
  NODE_OPTIONS+=("--no-trace_turbo_trimming")
fi

if [ -n "$TRACE_TURBO_JT" ] && [ $TRACE_TURBO_JT == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_jt")
else
  NODE_OPTIONS+=("--no-trace_turbo_jt")
fi

if [ -n "$TRACE_TURBO_CEQ" ] && [ $TRACE_TURBO_CEQ == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_ceq")
else
  NODE_OPTIONS+=("--no-trace_turbo_ceq")
fi

if [ -n "$TRACE_TURBO_LOOP" ] && [ $TRACE_TURBO_LOOP == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_loop")
else
  NODE_OPTIONS+=("--no-trace_turbo_loop")
fi

if [ -n "$TURBO_VERIFY" ] && [ $TURBO_VERIFY == "y" ]; then
  NODE_OPTIONS+=("--turbo_verify")
else
  NODE_OPTIONS+=("--no-turbo_verify")
fi

if [ -n "$TRACE_VERIFY_CSA" ] && [ $TRACE_VERIFY_CSA == "y" ]; then
  NODE_OPTIONS+=("--trace_verify_csa")
else
  NODE_OPTIONS+=("--no-trace_verify_csa")
fi

if [ -n "$TURBO_STATS" ] && [ $TURBO_STATS == "y" ]; then
  NODE_OPTIONS+=("--turbo_stats")
else
  NODE_OPTIONS+=("--no-turbo_stats")
fi

if [ -n "$TURBO_STATS_NVP" ] && [ $TURBO_STATS_NVP == "y" ]; then
  NODE_OPTIONS+=("--turbo_stats_nvp")
else
  NODE_OPTIONS+=("--no-turbo_stats_nvp")
fi

if [ -n "$TURBO_SPLITTING" ] && [ $TURBO_SPLITTING == "y" ]; then
  NODE_OPTIONS+=("--turbo_splitting")
else
  NODE_OPTIONS+=("--no-turbo_splitting")
fi

if [ -n "$FUNCTION_CONTEXT_SPECIALIZATION" ] && [ $FUNCTION_CONTEXT_SPECIALIZATION == "y" ]; then
  NODE_OPTIONS+=("--function_context_specialization")
else
  NODE_OPTIONS+=("--no-function_context_specialization")
fi

if [ -n "$TURBO_INLINING" ] && [ $TURBO_INLINING == "y" ]; then
  NODE_OPTIONS+=("--turbo_inlining")
else
  NODE_OPTIONS+=("--no-turbo_inlining")
fi

if [ -n "$TRACE_TURBO_INLINING" ] && [ $TRACE_TURBO_INLINING == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_inlining")
else
  NODE_OPTIONS+=("--no-trace_turbo_inlining")
fi

if [ -n "$TURBO_LOAD_ELIMINATION" ] && [ $TURBO_LOAD_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_load_elimination")
else
  NODE_OPTIONS+=("--no-turbo_load_elimination")
fi

if [ -n "$TRACE_TURBO_LOAD_ELIMINATION" ] && [ $TRACE_TURBO_LOAD_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--trace_turbo_load_elimination")
else
  NODE_OPTIONS+=("--no-trace_turbo_load_elimination")
fi

if [ -n "$TURBO_PROFILING" ] && [ $TURBO_PROFILING == "y" ]; then
  NODE_OPTIONS+=("--turbo_profiling")
else
  NODE_OPTIONS+=("--no-turbo_profiling")
fi

if [ -n "$TURBO_VERIFY_ALLOCATION" ] && [ $TURBO_VERIFY_ALLOCATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_verify_allocation")
else
  NODE_OPTIONS+=("--no-turbo_verify_allocation")
fi

if [ -n "$TURBO_MOVE_OPTIMIZATION" ] && [ $TURBO_MOVE_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_move_optimization")
else
  NODE_OPTIONS+=("--no-turbo_move_optimization")
fi

if [ -n "$TURBO_JT" ] && [ $TURBO_JT == "y" ]; then
  NODE_OPTIONS+=("--turbo_jt")
else
  NODE_OPTIONS+=("--no-turbo_jt")
fi

if [ -n "$TURBO_LOOP_PEELING" ] && [ $TURBO_LOOP_PEELING == "y" ]; then
  NODE_OPTIONS+=("--turbo_loop_peeling")
else
  NODE_OPTIONS+=("--no-turbo_loop_peeling")
fi

if [ -n "$TURBO_LOOP_VARIABLE" ] && [ $TURBO_LOOP_VARIABLE == "y" ]; then
  NODE_OPTIONS+=("--turbo_loop_variable")
else
  NODE_OPTIONS+=("--no-turbo_loop_variable")
fi

if [ -n "$TURBO_CF_OPTIMIZATION" ] && [ $TURBO_CF_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_cf_optimization")
else
  NODE_OPTIONS+=("--no-turbo_cf_optimization")
fi

if [ -n "$TURBO_FRAME_ELISION" ] && [ $TURBO_FRAME_ELISION == "y" ]; then
  NODE_OPTIONS+=("--turbo_frame_elision")
else
  NODE_OPTIONS+=("--no-turbo_frame_elision")
fi

if [ -n "$TURBO_ESCAPE" ] && [ $TURBO_ESCAPE == "y" ]; then
  NODE_OPTIONS+=("--turbo_escape")
else
  NODE_OPTIONS+=("--no-turbo_escape")
fi

if [ -n "$TURBO_INSTRUCTION_SCHEDULING" ] && [ $TURBO_INSTRUCTION_SCHEDULING == "y" ]; then
  NODE_OPTIONS+=("--turbo_instruction_scheduling")
else
  NODE_OPTIONS+=("--no-turbo_instruction_scheduling")
fi

if [ -n "$TURBO_STRESS_INSTRUCTION_SCHEDULING" ] && [ $TURBO_STRESS_INSTRUCTION_SCHEDULING == "y" ]; then
  NODE_OPTIONS+=("--turbo_stress_instruction_scheduling")
else
  NODE_OPTIONS+=("--no-turbo_stress_instruction_scheduling")
fi

if [ -n "$TURBO_STORE_ELIMINATION" ] && [ $TURBO_STORE_ELIMINATION == "y" ]; then
  NODE_OPTIONS+=("--turbo_store_elimination")
else
  NODE_OPTIONS+=("--no-turbo_store_elimination")
fi

if [ -n "$MINIMAL" ] && [ $MINIMAL == "y" ]; then
  NODE_OPTIONS+=("--minimal")
else
  NODE_OPTIONS+=("--no-minimal")
fi

if [ -n "$EXPOSE_WASM" ] && [ $EXPOSE_WASM == "y" ]; then
  NODE_OPTIONS+=("--expose_wasm")
else
  NODE_OPTIONS+=("--no-expose_wasm")
fi

if [ -n "$ASSUME_ASMJS_ORIGIN" ] && [ $ASSUME_ASMJS_ORIGIN == "y" ]; then
  NODE_OPTIONS+=("--assume_asmjs_origin")
else
  NODE_OPTIONS+=("--no-assume_asmjs_origin")
fi

if [ -n "$WASM_DISABLE_STRUCTURED_CLONING" ] && [ $WASM_DISABLE_STRUCTURED_CLONING == "y" ]; then
  NODE_OPTIONS+=("--wasm_disable_structured_cloning")
else
  NODE_OPTIONS+=("--no-wasm_disable_structured_cloning")
fi

if [ -n "$TRACE_WASM_DECODER" ] && [ $TRACE_WASM_DECODER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_decoder")
else
  NODE_OPTIONS+=("--no-trace_wasm_decoder")
fi

if [ -n "$TRACE_WASM_DECODE_TIME" ] && [ $TRACE_WASM_DECODE_TIME == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_decode_time")
else
  NODE_OPTIONS+=("--no-trace_wasm_decode_time")
fi

if [ -n "$TRACE_WASM_COMPILER" ] && [ $TRACE_WASM_COMPILER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_compiler")
else
  NODE_OPTIONS+=("--no-trace_wasm_compiler")
fi

if [ -n "$TRACE_WASM_INTERPRETER" ] && [ $TRACE_WASM_INTERPRETER == "y" ]; then
  NODE_OPTIONS+=("--trace_wasm_interpreter")
else
  NODE_OPTIONS+=("--no-trace_wasm_interpreter")
fi

if [ -n "$WASM_BREAK_ON_DECODER_ERROR" ] && [ $WASM_BREAK_ON_DECODER_ERROR == "y" ]; then
  NODE_OPTIONS+=("--wasm_break_on_decoder_error")
else
  NODE_OPTIONS+=("--no-wasm_break_on_decoder_error")
fi

if [ -n "$VALIDATE_ASM" ] && [ $VALIDATE_ASM == "y" ]; then
  NODE_OPTIONS+=("--validate_asm")
else
  NODE_OPTIONS+=("--no-validate_asm")
fi

if [ -n "$SUPPRESS_ASM_MESSAGES" ] && [ $SUPPRESS_ASM_MESSAGES == "y" ]; then
  NODE_OPTIONS+=("--suppress_asm_messages")
else
  NODE_OPTIONS+=("--no-suppress_asm_messages")
fi

if [ -n "$TRACE_ASM_TIME" ] && [ $TRACE_ASM_TIME == "y" ]; then
  NODE_OPTIONS+=("--trace_asm_time")
else
  NODE_OPTIONS+=("--no-trace_asm_time")
fi

if [ -n "$DUMP_WASM_MODULE" ] && [ $DUMP_WASM_MODULE == "y" ]; then
  NODE_OPTIONS+=("--dump_wasm_module")
else
  NODE_OPTIONS+=("--no-dump_wasm_module")
fi

if [ -n "$WASM_OPT" ] && [ $WASM_OPT == "y" ]; then
  NODE_OPTIONS+=("--wasm_opt")
else
  NODE_OPTIONS+=("--no-wasm_opt")
fi

if [ -n "$WASM_NO_BOUNDS_CHECKS" ] && [ $WASM_NO_BOUNDS_CHECKS == "y" ]; then
  NODE_OPTIONS+=("--wasm_no_bounds_checks")
else
  NODE_OPTIONS+=("--no-wasm_no_bounds_checks")
fi

if [ -n "$WASM_NO_STACK_CHECKS" ] && [ $WASM_NO_STACK_CHECKS == "y" ]; then
  NODE_OPTIONS+=("--wasm_no_stack_checks")
else
  NODE_OPTIONS+=("--no-wasm_no_stack_checks")
fi

if [ -n "$WASM_TRAP_HANDLER" ] && [ $WASM_TRAP_HANDLER == "y" ]; then
  NODE_OPTIONS+=("--wasm_trap_handler")
else
  NODE_OPTIONS+=("--no-wasm_trap_handler")
fi

if [ -n "$WASM_GUARD_PAGES" ] && [ $WASM_GUARD_PAGES == "y" ]; then
  NODE_OPTIONS+=("--wasm_guard_pages")
else
  NODE_OPTIONS+=("--no-wasm_guard_pages")
fi

if [ -n "$WASM_CODE_FUZZER_GEN_TEST" ] && [ $WASM_CODE_FUZZER_GEN_TEST == "y" ]; then
  NODE_OPTIONS+=("--wasm_code_fuzzer_gen_test")
else
  NODE_OPTIONS+=("--no-wasm_code_fuzzer_gen_test")
fi

if [ -n "$PRINT_WASM_CODE" ] && [ $PRINT_WASM_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_wasm_code")
else
  NODE_OPTIONS+=("--no-print_wasm_code")
fi

if [ -n "$TRACE_OPT_VERBOSE" ] && [ $TRACE_OPT_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_opt_verbose")
else
  NODE_OPTIONS+=("--no-trace_opt_verbose")
fi

if [ -n "$EXPERIMENTAL_NEW_SPACE_GROWTH_HEURISTIC" ] && [ $EXPERIMENTAL_NEW_SPACE_GROWTH_HEURISTIC == "y" ]; then
  NODE_OPTIONS+=("--experimental_new_space_growth_heuristic")
else
  NODE_OPTIONS+=("--no-experimental_new_space_growth_heuristic")
fi

if [ -n "$GC_GLOBAL" ] && [ $GC_GLOBAL == "y" ]; then
  NODE_OPTIONS+=("--gc_global")
else
  NODE_OPTIONS+=("--no-gc_global")
fi

if [ -n "$TRACE_GC" ] && [ $TRACE_GC == "y" ]; then
  NODE_OPTIONS+=("--trace_gc")
else
  NODE_OPTIONS+=("--no-trace_gc")
fi

if [ -n "$TRACE_GC_NVP" ] && [ $TRACE_GC_NVP == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_nvp")
else
  NODE_OPTIONS+=("--no-trace_gc_nvp")
fi

if [ -n "$TRACE_GC_IGNORE_SCAVENGER" ] && [ $TRACE_GC_IGNORE_SCAVENGER == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_ignore_scavenger")
else
  NODE_OPTIONS+=("--no-trace_gc_ignore_scavenger")
fi

if [ -n "$TRACE_IDLE_NOTIFICATION" ] && [ $TRACE_IDLE_NOTIFICATION == "y" ]; then
  NODE_OPTIONS+=("--trace_idle_notification")
else
  NODE_OPTIONS+=("--no-trace_idle_notification")
fi

if [ -n "$TRACE_IDLE_NOTIFICATION_VERBOSE" ] && [ $TRACE_IDLE_NOTIFICATION_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_idle_notification_verbose")
else
  NODE_OPTIONS+=("--no-trace_idle_notification_verbose")
fi

if [ -n "$TRACE_GC_VERBOSE" ] && [ $TRACE_GC_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_verbose")
else
  NODE_OPTIONS+=("--no-trace_gc_verbose")
fi

if [ -n "$TRACE_FRAGMENTATION" ] && [ $TRACE_FRAGMENTATION == "y" ]; then
  NODE_OPTIONS+=("--trace_fragmentation")
else
  NODE_OPTIONS+=("--no-trace_fragmentation")
fi

if [ -n "$TRACE_FRAGMENTATION_VERBOSE" ] && [ $TRACE_FRAGMENTATION_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--trace_fragmentation_verbose")
else
  NODE_OPTIONS+=("--no-trace_fragmentation_verbose")
fi

if [ -n "$TRACE_EVACUATION" ] && [ $TRACE_EVACUATION == "y" ]; then
  NODE_OPTIONS+=("--trace_evacuation")
else
  NODE_OPTIONS+=("--no-trace_evacuation")
fi

if [ -n "$TRACE_MUTATOR_UTILIZATION" ] && [ $TRACE_MUTATOR_UTILIZATION == "y" ]; then
  NODE_OPTIONS+=("--trace_mutator_utilization")
else
  NODE_OPTIONS+=("--no-trace_mutator_utilization")
fi

if [ -n "$INCREMENTAL_MARKING" ] && [ $INCREMENTAL_MARKING == "y" ]; then
  NODE_OPTIONS+=("--incremental_marking")
else
  NODE_OPTIONS+=("--no-incremental_marking")
fi

if [ -n "$INCREMENTAL_MARKING_WRAPPERS" ] && [ $INCREMENTAL_MARKING_WRAPPERS == "y" ]; then
  NODE_OPTIONS+=("--incremental_marking_wrappers")
else
  NODE_OPTIONS+=("--no-incremental_marking_wrappers")
fi

if [ -n "$MINOR_MC" ] && [ $MINOR_MC == "y" ]; then
  NODE_OPTIONS+=("--minor_mc")
else
  NODE_OPTIONS+=("--no-minor_mc")
fi

if [ -n "$BLACK_ALLOCATION" ] && [ $BLACK_ALLOCATION == "y" ]; then
  NODE_OPTIONS+=("--black_allocation")
else
  NODE_OPTIONS+=("--no-black_allocation")
fi

if [ -n "$CONCURRENT_SWEEPING" ] && [ $CONCURRENT_SWEEPING == "y" ]; then
  NODE_OPTIONS+=("--concurrent_sweeping")
else
  NODE_OPTIONS+=("--no-concurrent_sweeping")
fi

if [ -n "$PARALLEL_COMPACTION" ] && [ $PARALLEL_COMPACTION == "y" ]; then
  NODE_OPTIONS+=("--parallel_compaction")
else
  NODE_OPTIONS+=("--no-parallel_compaction")
fi

if [ -n "$PARALLEL_POINTER_UPDATE" ] && [ $PARALLEL_POINTER_UPDATE == "y" ]; then
  NODE_OPTIONS+=("--parallel_pointer_update")
else
  NODE_OPTIONS+=("--no-parallel_pointer_update")
fi

if [ -n "$TRACE_INCREMENTAL_MARKING" ] && [ $TRACE_INCREMENTAL_MARKING == "y" ]; then
  NODE_OPTIONS+=("--trace_incremental_marking")
else
  NODE_OPTIONS+=("--no-trace_incremental_marking")
fi

if [ -n "$TRACK_GC_OBJECT_STATS" ] && [ $TRACK_GC_OBJECT_STATS == "y" ]; then
  NODE_OPTIONS+=("--track_gc_object_stats")
else
  NODE_OPTIONS+=("--no-track_gc_object_stats")
fi

if [ -n "$TRACE_GC_OBJECT_STATS" ] && [ $TRACE_GC_OBJECT_STATS == "y" ]; then
  NODE_OPTIONS+=("--trace_gc_object_stats")
else
  NODE_OPTIONS+=("--no-trace_gc_object_stats")
fi

if [ -n "$TRACK_DETACHED_CONTEXTS" ] && [ $TRACK_DETACHED_CONTEXTS == "y" ]; then
  NODE_OPTIONS+=("--track_detached_contexts")
else
  NODE_OPTIONS+=("--no-track_detached_contexts")
fi

if [ -n "$TRACE_DETACHED_CONTEXTS" ] && [ $TRACE_DETACHED_CONTEXTS == "y" ]; then
  NODE_OPTIONS+=("--trace_detached_contexts")
else
  NODE_OPTIONS+=("--no-trace_detached_contexts")
fi

if [ -n "$MOVE_OBJECT_START" ] && [ $MOVE_OBJECT_START == "y" ]; then
  NODE_OPTIONS+=("--move_object_start")
else
  NODE_OPTIONS+=("--no-move_object_start")
fi

if [ -n "$MEMORY_REDUCER" ] && [ $MEMORY_REDUCER == "y" ]; then
  NODE_OPTIONS+=("--memory_reducer")
else
  NODE_OPTIONS+=("--no-memory_reducer")
fi

if [ -n "$ALWAYS_COMPACT" ] && [ $ALWAYS_COMPACT == "y" ]; then
  NODE_OPTIONS+=("--always_compact")
else
  NODE_OPTIONS+=("--no-always_compact")
fi

if [ -n "$NEVER_COMPACT" ] && [ $NEVER_COMPACT == "y" ]; then
  NODE_OPTIONS+=("--never_compact")
else
  NODE_OPTIONS+=("--no-never_compact")
fi

if [ -n "$COMPACT_CODE_SPACE" ] && [ $COMPACT_CODE_SPACE == "y" ]; then
  NODE_OPTIONS+=("--compact_code_space")
else
  NODE_OPTIONS+=("--no-compact_code_space")
fi

if [ -n "$CLEANUP_CODE_CACHES_AT_GC" ] && [ $CLEANUP_CODE_CACHES_AT_GC == "y" ]; then
  NODE_OPTIONS+=("--cleanup_code_caches_at_gc")
else
  NODE_OPTIONS+=("--no-cleanup_code_caches_at_gc")
fi

if [ -n "$USE_MARKING_PROGRESS_BAR" ] && [ $USE_MARKING_PROGRESS_BAR == "y" ]; then
  NODE_OPTIONS+=("--use_marking_progress_bar")
else
  NODE_OPTIONS+=("--no-use_marking_progress_bar")
fi

if [ -n "$FORCE_MARKING_DEQUE_OVERFLOWS" ] && [ $FORCE_MARKING_DEQUE_OVERFLOWS == "y" ]; then
  NODE_OPTIONS+=("--force_marking_deque_overflows")
else
  NODE_OPTIONS+=("--no-force_marking_deque_overflows")
fi

if [ -n "$STRESS_COMPACTION" ] && [ $STRESS_COMPACTION == "y" ]; then
  NODE_OPTIONS+=("--stress_compaction")
else
  NODE_OPTIONS+=("--no-stress_compaction")
fi

if [ -n "$MANUAL_EVACUATION_CANDIDATES_SELECTION" ] && [ $MANUAL_EVACUATION_CANDIDATES_SELECTION == "y" ]; then
  NODE_OPTIONS+=("--manual_evacuation_candidates_selection")
else
  NODE_OPTIONS+=("--no-manual_evacuation_candidates_selection")
fi

if [ -n "$FAST_PROMOTION_NEW_SPACE" ] && [ $FAST_PROMOTION_NEW_SPACE == "y" ]; then
  NODE_OPTIONS+=("--fast_promotion_new_space")
else
  NODE_OPTIONS+=("--no-fast_promotion_new_space")
fi

if [ -n "$DEBUG_CODE" ] && [ $DEBUG_CODE == "y" ]; then
  NODE_OPTIONS+=("--debug_code")
else
  NODE_OPTIONS+=("--no-debug_code")
fi

if [ -n "$CODE_COMMENTS" ] && [ $CODE_COMMENTS == "y" ]; then
  NODE_OPTIONS+=("--code_comments")
else
  NODE_OPTIONS+=("--no-code_comments")
fi

if [ -n "$ENABLE_SSE3" ] && [ $ENABLE_SSE3 == "y" ]; then
  NODE_OPTIONS+=("--enable_sse3")
else
  NODE_OPTIONS+=("--no-enable_sse3")
fi

if [ -n "$ENABLE_SSSE" ] && [ $ENABLE_SSSE == "y" ]; then
  NODE_OPTIONS+=("--enable_ssse3")
else
  NODE_OPTIONS+=("--no-enable_ssse3")
fi

if [ -n "$ENABLE_SSE4_1" ] && [ $ENABLE_SSE4_1 == "y" ]; then
  NODE_OPTIONS+=("--enable_sse4_1")
else
  NODE_OPTIONS+=("--no-enable_sse4_1")
fi

if [ -n "$ENABLE_SAHF" ] && [ $ENABLE_SAHF == "y" ]; then
  NODE_OPTIONS+=("--enable_sahf")
else
  NODE_OPTIONS+=("--no-enable_sahf")
fi

if [ -n "$ENABLE_AVX" ] && [ $ENABLE_AVX == "y" ]; then
  NODE_OPTIONS+=("--enable_avx")
else
  NODE_OPTIONS+=("--no-enable_avx")
fi

if [ -n "$ENABLE_FMA" ] && [ $ENABLE_FMA == "y" ]; then
  NODE_OPTIONS+=("--enable_fma3")
else
  NODE_OPTIONS+=("--no-enable_fma3")
fi

if [ -n "$ENABLE_BMI1" ] && [ $ENABLE_BMI1 == "y" ]; then
  NODE_OPTIONS+=("--enable_bmi1")
else
  NODE_OPTIONS+=("--no-enable_bmi1")
fi

if [ -n "$ENABLE_BMI2" ] && [ $ENABLE_BMI2 == "y" ]; then
  NODE_OPTIONS+=("--enable_bmi2")
else
  NODE_OPTIONS+=("--no-enable_bmi2")
fi

if [ -n "$ENABLE_LZCNT" ] && [ $ENABLE_LZCNT == "y" ]; then
  NODE_OPTIONS+=("--enable_lzcnt")
else
  NODE_OPTIONS+=("--no-enable_lzcnt")
fi

if [ -n "$ENABLE_POPCNT" ] && [ $ENABLE_POPCNT == "y" ]; then
  NODE_OPTIONS+=("--enable_popcnt")
else
  NODE_OPTIONS+=("--no-enable_popcnt")
fi

if [ -n "$ENABLE_VLDR_IMM" ] && [ $ENABLE_VLDR_IMM == "y" ]; then
  NODE_OPTIONS+=("--enable_vldr_imm")
else
  NODE_OPTIONS+=("--no-enable_vldr_imm")
fi

if [ -n "$FORCE_LONG_BRANCHES" ] && [ $FORCE_LONG_BRANCHES == "y" ]; then
  NODE_OPTIONS+=("--force_long_branches")
else
  NODE_OPTIONS+=("--no-force_long_branches")
fi

if [ -n "$ENABLE_ARMV7" ] && [ $ENABLE_ARMV7 == "y" ]; then
  NODE_OPTIONS+=("--enable_armv7")
else
  NODE_OPTIONS+=("--no-enable_armv7")
fi

if [ -n "$ENABLE_VFP" ] && [ $ENABLE_VFP == "y" ]; then
  NODE_OPTIONS+=("--enable_vfp3")
else
  NODE_OPTIONS+=("--no-enable_vfp3")
fi

if [ -n "$ENABLE_" ] && [ $ENABLE_ == "y" ]; then
  NODE_OPTIONS+=("--enable_32dregs")
else
  NODE_OPTIONS+=("--no-enable_32dregs")
fi

if [ -n "$ENABLE_NEON" ] && [ $ENABLE_NEON == "y" ]; then
  NODE_OPTIONS+=("--enable_neon")
else
  NODE_OPTIONS+=("--no-enable_neon")
fi

if [ -n "$ENABLE_SUDIV" ] && [ $ENABLE_SUDIV == "y" ]; then
  NODE_OPTIONS+=("--enable_sudiv")
else
  NODE_OPTIONS+=("--no-enable_sudiv")
fi

if [ -n "$ENABLE_ARMV8" ] && [ $ENABLE_ARMV8 == "y" ]; then
  NODE_OPTIONS+=("--enable_armv8")
else
  NODE_OPTIONS+=("--no-enable_armv8")
fi

if [ -n "$ENABLE_REGEXP_UNALIGNED_ACCESSES" ] && [ $ENABLE_REGEXP_UNALIGNED_ACCESSES == "y" ]; then
  NODE_OPTIONS+=("--enable_regexp_unaligned_accesses")
else
  NODE_OPTIONS+=("--no-enable_regexp_unaligned_accesses")
fi

if [ -n "$SCRIPT_STREAMING" ] && [ $SCRIPT_STREAMING == "y" ]; then
  NODE_OPTIONS+=("--script_streaming")
else
  NODE_OPTIONS+=("--no-script_streaming")
fi

if [ -n "$DISABLE_OLD_API_ACCESSORS" ] && [ $DISABLE_OLD_API_ACCESSORS == "y" ]; then
  NODE_OPTIONS+=("--disable_old_api_accessors")
else
  NODE_OPTIONS+=("--no-disable_old_api_accessors")
fi

if [ -n "$EXPOSE_FREE_BUFFER" ] && [ $EXPOSE_FREE_BUFFER == "y" ]; then
  NODE_OPTIONS+=("--expose_free_buffer")
else
  NODE_OPTIONS+=("--no-expose_free_buffer")
fi

if [ -n "$EXPOSE_GC" ] && [ $EXPOSE_GC == "y" ]; then
  NODE_OPTIONS+=("--expose_gc")
else
  NODE_OPTIONS+=("--no-expose_gc")
fi

if [ -n "$EXPOSE_EXTERNALIZE_STRING" ] && [ $EXPOSE_EXTERNALIZE_STRING == "y" ]; then
  NODE_OPTIONS+=("--expose_externalize_string")
else
  NODE_OPTIONS+=("--no-expose_externalize_string")
fi

if [ -n "$EXPOSE_TRIGGER_FAILURE" ] && [ $EXPOSE_TRIGGER_FAILURE == "y" ]; then
  NODE_OPTIONS+=("--expose_trigger_failure")
else
  NODE_OPTIONS+=("--no-expose_trigger_failure")
fi

if [ -n "$BUILTINS_IN_STACK_TRACES" ] && [ $BUILTINS_IN_STACK_TRACES == "y" ]; then
  NODE_OPTIONS+=("--builtins_in_stack_traces")
else
  NODE_OPTIONS+=("--no-builtins_in_stack_traces")
fi

if [ -n "$ALLOW_UNSAFE_FUNCTION_CONSTRUCTOR" ] && [ $ALLOW_UNSAFE_FUNCTION_CONSTRUCTOR == "y" ]; then
  NODE_OPTIONS+=("--allow_unsafe_function_constructor")
else
  NODE_OPTIONS+=("--no-allow_unsafe_function_constructor")
fi

if [ -n "$INLINE_NEW" ] && [ $INLINE_NEW == "y" ]; then
  NODE_OPTIONS+=("--inline_new")
else
  NODE_OPTIONS+=("--no-inline_new")
fi

if [ -n "$TRACE_CODEGEN" ] && [ $TRACE_CODEGEN == "y" ]; then
  NODE_OPTIONS+=("--trace_codegen")
else
  NODE_OPTIONS+=("--no-trace_codegen")
fi

if [ -n "$LAZY" ] && [ $LAZY == "y" ]; then
  NODE_OPTIONS+=("--lazy")
else
  NODE_OPTIONS+=("--no-lazy")
fi

if [ -n "$TRACE_OPT" ] && [ $TRACE_OPT == "y" ]; then
  NODE_OPTIONS+=("--trace_opt")
else
  NODE_OPTIONS+=("--no-trace_opt")
fi

if [ -n "$TRACE_OPT_STATS" ] && [ $TRACE_OPT_STATS == "y" ]; then
  NODE_OPTIONS+=("--trace_opt_stats")
else
  NODE_OPTIONS+=("--no-trace_opt_stats")
fi

if [ -n "$TRACE_FILE_NAMES" ] && [ $TRACE_FILE_NAMES == "y" ]; then
  NODE_OPTIONS+=("--trace_file_names")
else
  NODE_OPTIONS+=("--no-trace_file_names")
fi

if [ -n "$OPT" ] && [ $OPT == "y" ]; then
  NODE_OPTIONS+=("--opt")
else
  NODE_OPTIONS+=("--no-opt")
fi

if [ -n "$ALWAYS_OPT" ] && [ $ALWAYS_OPT == "y" ]; then
  NODE_OPTIONS+=("--always_opt")
else
  NODE_OPTIONS+=("--no-always_opt")
fi

if [ -n "$ALWAYS_OSR" ] && [ $ALWAYS_OSR == "y" ]; then
  NODE_OPTIONS+=("--always_osr")
else
  NODE_OPTIONS+=("--no-always_osr")
fi

if [ -n "$PREPARE_ALWAYS_OPT" ] && [ $PREPARE_ALWAYS_OPT == "y" ]; then
  NODE_OPTIONS+=("--prepare_always_opt")
else
  NODE_OPTIONS+=("--no-prepare_always_opt")
fi

if [ -n "$TRACE_DEOPT" ] && [ $TRACE_DEOPT == "y" ]; then
  NODE_OPTIONS+=("--trace_deopt")
else
  NODE_OPTIONS+=("--no-trace_deopt")
fi

if [ -n "$SERIALIZE_TOPLEVEL" ] && [ $SERIALIZE_TOPLEVEL == "y" ]; then
  NODE_OPTIONS+=("--serialize_toplevel")
else
  NODE_OPTIONS+=("--no-serialize_toplevel")
fi

if [ -n "$SERIALIZE_EAGER" ] && [ $SERIALIZE_EAGER == "y" ]; then
  NODE_OPTIONS+=("--serialize_eager")
else
  NODE_OPTIONS+=("--no-serialize_eager")
fi

if [ -n "$TRACE_SERIALIZER" ] && [ $TRACE_SERIALIZER == "y" ]; then
  NODE_OPTIONS+=("--trace_serializer")
else
  NODE_OPTIONS+=("--no-trace_serializer")
fi

if [ -n "$COMPILATION_CACHE" ] && [ $COMPILATION_CACHE == "y" ]; then
  NODE_OPTIONS+=("--compilation_cache")
else
  NODE_OPTIONS+=("--no-compilation_cache")
fi

if [ -n "$CACHE_PROTOTYPE_TRANSITIONS" ] && [ $CACHE_PROTOTYPE_TRANSITIONS == "y" ]; then
  NODE_OPTIONS+=("--cache_prototype_transitions")
else
  NODE_OPTIONS+=("--no-cache_prototype_transitions")
fi

if [ -n "$COMPILER_DISPATCHER" ] && [ $COMPILER_DISPATCHER == "y" ]; then
  NODE_OPTIONS+=("--compiler_dispatcher")
else
  NODE_OPTIONS+=("--no-compiler_dispatcher")
fi

if [ -n "$TRACE_COMPILER_DISPATCHER" ] && [ $TRACE_COMPILER_DISPATCHER == "y" ]; then
  NODE_OPTIONS+=("--trace_compiler_dispatcher")
else
  NODE_OPTIONS+=("--no-trace_compiler_dispatcher")
fi

if [ -n "$TRACE_COMPILER_DISPATCHER_JOBS" ] && [ $TRACE_COMPILER_DISPATCHER_JOBS == "y" ]; then
  NODE_OPTIONS+=("--trace_compiler_dispatcher_jobs")
else
  NODE_OPTIONS+=("--no-trace_compiler_dispatcher_jobs")
fi

if [ -n "$TRACE_JS_ARRAY_ABUSE" ] && [ $TRACE_JS_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_js_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_js_array_abuse")
fi

if [ -n "$TRACE_EXTERNAL_ARRAY_ABUSE" ] && [ $TRACE_EXTERNAL_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_external_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_external_array_abuse")
fi

if [ -n "$TRACE_ARRAY_ABUSE" ] && [ $TRACE_ARRAY_ABUSE == "y" ]; then
  NODE_OPTIONS+=("--trace_array_abuse")
else
  NODE_OPTIONS+=("--no-trace_array_abuse")
fi

if [ -n "$ENABLE_LIVEEDIT" ] && [ $ENABLE_LIVEEDIT == "y" ]; then
  NODE_OPTIONS+=("--enable_liveedit")
else
  NODE_OPTIONS+=("--no-enable_liveedit")
fi

if [ -n "$TRACE_SIDE_EFFECT_FREE_DEBUG_EVALUATE" ] && [ $TRACE_SIDE_EFFECT_FREE_DEBUG_EVALUATE == "y" ]; then
  NODE_OPTIONS+=("--trace_side_effect_free_debug_evaluate")
else
  NODE_OPTIONS+=("--no-trace_side_effect_free_debug_evaluate")
fi

if [ -n "$HARD_ABORT" ] && [ $HARD_ABORT == "y" ]; then
  NODE_OPTIONS+=("--hard_abort")
else
  NODE_OPTIONS+=("--no-hard_abort")
fi

if [ -n "$CLEAR_EXCEPTIONS_ON_JS_ENTRY" ] && [ $CLEAR_EXCEPTIONS_ON_JS_ENTRY == "y" ]; then
  NODE_OPTIONS+=("--clear_exceptions_on_js_entry")
else
  NODE_OPTIONS+=("--no-clear_exceptions_on_js_entry")
fi

if [ -n "$HEAP_PROFILER_TRACE_OBJECTS" ] && [ $HEAP_PROFILER_TRACE_OBJECTS == "y" ]; then
  NODE_OPTIONS+=("--heap_profiler_trace_objects")
else
  NODE_OPTIONS+=("--no-heap_profiler_trace_objects")
fi

if [ -n "$SAMPLING_HEAP_PROFILER_SUPPRESS_RANDOMNESS" ] && [ $SAMPLING_HEAP_PROFILER_SUPPRESS_RANDOMNESS == "y" ]; then
  NODE_OPTIONS+=("--sampling_heap_profiler_suppress_randomness")
else
  NODE_OPTIONS+=("--no-sampling_heap_profiler_suppress_randomness")
fi

if [ -n "$USE_IDLE_NOTIFICATION" ] && [ $USE_IDLE_NOTIFICATION == "y" ]; then
  NODE_OPTIONS+=("--use_idle_notification")
else
  NODE_OPTIONS+=("--no-use_idle_notification")
fi

if [ -n "$USE_IC" ] && [ $USE_IC == "y" ]; then
  NODE_OPTIONS+=("--use_ic")
else
  NODE_OPTIONS+=("--no-use_ic")
fi

if [ -n "$TRACE_IC" ] && [ $TRACE_IC == "y" ]; then
  NODE_OPTIONS+=("--trace_ic")
else
  NODE_OPTIONS+=("--no-trace_ic")
fi

if [ -n "$NATIVE_CODE_COUNTERS" ] && [ $NATIVE_CODE_COUNTERS == "y" ]; then
  NODE_OPTIONS+=("--native_code_counters")
else
  NODE_OPTIONS+=("--no-native_code_counters")
fi

if [ -n "$THIN_STRINGS" ] && [ $THIN_STRINGS == "y" ]; then
  NODE_OPTIONS+=("--thin_strings")
else
  NODE_OPTIONS+=("--no-thin_strings")
fi

if [ -n "$TRACE_WEAK_ARRAYS" ] && [ $TRACE_WEAK_ARRAYS == "y" ]; then
  NODE_OPTIONS+=("--trace_weak_arrays")
else
  NODE_OPTIONS+=("--no-trace_weak_arrays")
fi

if [ -n "$TRACE_PROTOTYPE_USERS" ] && [ $TRACE_PROTOTYPE_USERS == "y" ]; then
  NODE_OPTIONS+=("--trace_prototype_users")
else
  NODE_OPTIONS+=("--no-trace_prototype_users")
fi

if [ -n "$USE_VERBOSE_PRINTER" ] && [ $USE_VERBOSE_PRINTER == "y" ]; then
  NODE_OPTIONS+=("--use_verbose_printer")
else
  NODE_OPTIONS+=("--no-use_verbose_printer")
fi

if [ -n "$TRACE_FOR_IN_ENUMERATE" ] && [ $TRACE_FOR_IN_ENUMERATE == "y" ]; then
  NODE_OPTIONS+=("--trace_for_in_enumerate")
else
  NODE_OPTIONS+=("--no-trace_for_in_enumerate")
fi

if [ -n "$ALLOW_NATIVES_SYNTAX" ] && [ $ALLOW_NATIVES_SYNTAX == "y" ]; then
  NODE_OPTIONS+=("--allow_natives_syntax")
else
  NODE_OPTIONS+=("--no-allow_natives_syntax")
fi

if [ -n "$TRACE_PARSE" ] && [ $TRACE_PARSE == "y" ]; then
  NODE_OPTIONS+=("--trace_parse")
else
  NODE_OPTIONS+=("--no-trace_parse")
fi

if [ -n "$TRACE_PREPARSE" ] && [ $TRACE_PREPARSE == "y" ]; then
  NODE_OPTIONS+=("--trace_preparse")
else
  NODE_OPTIONS+=("--no-trace_preparse")
fi

if [ -n "$LAZY_INNER_FUNCTIONS" ] && [ $LAZY_INNER_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--lazy_inner_functions")
else
  NODE_OPTIONS+=("--no-lazy_inner_functions")
fi

if [ -n "$AGGRESSIVE_LAZY_INNER_FUNCTIONS" ] && [ $AGGRESSIVE_LAZY_INNER_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--aggressive_lazy_inner_functions")
else
  NODE_OPTIONS+=("--no-aggressive_lazy_inner_functions")
fi

if [ -n "$PREPARSER_SCOPE_ANALYSIS" ] && [ $PREPARSER_SCOPE_ANALYSIS == "y" ]; then
  NODE_OPTIONS+=("--preparser_scope_analysis")
else
  NODE_OPTIONS+=("--no-preparser_scope_analysis")
fi

if [ -n "$TRACE_SIM" ] && [ $TRACE_SIM == "y" ]; then
  NODE_OPTIONS+=("--trace_sim")
else
  NODE_OPTIONS+=("--no-trace_sim")
fi

if [ -n "$DEBUG_SIM" ] && [ $DEBUG_SIM == "y" ]; then
  NODE_OPTIONS+=("--debug_sim")
else
  NODE_OPTIONS+=("--no-debug_sim")
fi

if [ -n "$CHECK_ICACHE" ] && [ $CHECK_ICACHE == "y" ]; then
  NODE_OPTIONS+=("--check_icache")
else
  NODE_OPTIONS+=("--no-check_icache")
fi

if [ -n "$LOG_REGS_MODIFIED" ] && [ $LOG_REGS_MODIFIED == "y" ]; then
  NODE_OPTIONS+=("--log_regs_modified")
else
  NODE_OPTIONS+=("--no-log_regs_modified")
fi

if [ -n "$LOG_COLOUR" ] && [ $LOG_COLOUR == "y" ]; then
  NODE_OPTIONS+=("--log_colour")
else
  NODE_OPTIONS+=("--no-log_colour")
fi

if [ -n "$IGNORE_ASM_UNIMPLEMENTED_BREAK" ] && [ $IGNORE_ASM_UNIMPLEMENTED_BREAK == "y" ]; then
  NODE_OPTIONS+=("--ignore_asm_unimplemented_break")
else
  NODE_OPTIONS+=("--no-ignore_asm_unimplemented_break")
fi

if [ -n "$TRACE_SIM_MESSAGES" ] && [ $TRACE_SIM_MESSAGES == "y" ]; then
  NODE_OPTIONS+=("--trace_sim_messages")
else
  NODE_OPTIONS+=("--no-trace_sim_messages")
fi

if [ -n "$STACK_TRACE_ON_ILLEGAL" ] && [ $STACK_TRACE_ON_ILLEGAL == "y" ]; then
  NODE_OPTIONS+=("--stack_trace_on_illegal")
else
  NODE_OPTIONS+=("--no-stack_trace_on_illegal")
fi

if [ -n "$ABORT_ON_UNCAUGHT_EXCEPTION" ] && [ $ABORT_ON_UNCAUGHT_EXCEPTION == "y" ]; then
  NODE_OPTIONS+=("--abort_on_uncaught_exception")
else
  NODE_OPTIONS+=("--no-abort_on_uncaught_exception")
fi

if [ -n "$RANDOMIZE_HASHES" ] && [ $RANDOMIZE_HASHES == "y" ]; then
  NODE_OPTIONS+=("--randomize_hashes")
else
  NODE_OPTIONS+=("--no-randomize_hashes")
fi

if [ -n "$TRACE_RAIL" ] && [ $TRACE_RAIL == "y" ]; then
  NODE_OPTIONS+=("--trace_rail")
else
  NODE_OPTIONS+=("--no-trace_rail")
fi

if [ -n "$PRINT_ALL_EXCEPTIONS" ] && [ $PRINT_ALL_EXCEPTIONS == "y" ]; then
  NODE_OPTIONS+=("--print_all_exceptions")
else
  NODE_OPTIONS+=("--no-print_all_exceptions")
fi

if [ -n "$RUNTIME_CALL_STATS" ] && [ $RUNTIME_CALL_STATS == "y" ]; then
  NODE_OPTIONS+=("--runtime_call_stats")
else
  NODE_OPTIONS+=("--no-runtime_call_stats")
fi

if [ -n "$PROFILE_DESERIALIZATION" ] && [ $PROFILE_DESERIALIZATION == "y" ]; then
  NODE_OPTIONS+=("--profile_deserialization")
else
  NODE_OPTIONS+=("--no-profile_deserialization")
fi

if [ -n "$SERIALIZATION_STATISTICS" ] && [ $SERIALIZATION_STATISTICS == "y" ]; then
  NODE_OPTIONS+=("--serialization_statistics")
else
  NODE_OPTIONS+=("--no-serialization_statistics")
fi

if [ -n "$REGEXP_OPTIMIZATION" ] && [ $REGEXP_OPTIMIZATION == "y" ]; then
  NODE_OPTIONS+=("--regexp_optimization")
else
  NODE_OPTIONS+=("--no-regexp_optimization")
fi

if [ -n "$TESTING_BOOL_FLAG" ] && [ $TESTING_BOOL_FLAG == "y" ]; then
  NODE_OPTIONS+=("--testing_bool_flag")
else
  NODE_OPTIONS+=("--no-testing_bool_flag")
fi

if [ -n "$TESTING_MAYBE_BOOL_FLAG" ] && [ $TESTING_MAYBE_BOOL_FLAG == "y" ]; then
  NODE_OPTIONS+=("--testing_maybe_bool_flag")
else
  NODE_OPTIONS+=("--no-testing_maybe_bool_flag")
fi

if [ -n "$DUMP_COUNTERS" ] && [ $DUMP_COUNTERS == "y" ]; then
  NODE_OPTIONS+=("--dump_counters")
else
  NODE_OPTIONS+=("--no-dump_counters")
fi

if [ -n "$DUMP_COUNTERS_NVP" ] && [ $DUMP_COUNTERS_NVP == "y" ]; then
  NODE_OPTIONS+=("--dump_counters_nvp")
else
  NODE_OPTIONS+=("--no-dump_counters_nvp")
fi

if [ -n "$LOG" ] && [ $LOG == "y" ]; then
  NODE_OPTIONS+=("--log")
else
  NODE_OPTIONS+=("--no-log")
fi

if [ -n "$LOG_ALL" ] && [ $LOG_ALL == "y" ]; then
  NODE_OPTIONS+=("--log_all")
else
  NODE_OPTIONS+=("--no-log_all")
fi

if [ -n "$LOG_API" ] && [ $LOG_API == "y" ]; then
  NODE_OPTIONS+=("--log_api")
else
  NODE_OPTIONS+=("--no-log_api")
fi

if [ -n "$LOG_CODE" ] && [ $LOG_CODE == "y" ]; then
  NODE_OPTIONS+=("--log_code")
else
  NODE_OPTIONS+=("--no-log_code")
fi

if [ -n "$LOG_GC" ] && [ $LOG_GC == "y" ]; then
  NODE_OPTIONS+=("--log_gc")
else
  NODE_OPTIONS+=("--no-log_gc")
fi

if [ -n "$LOG_HANDLES" ] && [ $LOG_HANDLES == "y" ]; then
  NODE_OPTIONS+=("--log_handles")
else
  NODE_OPTIONS+=("--no-log_handles")
fi

if [ -n "$LOG_SUSPECT" ] && [ $LOG_SUSPECT == "y" ]; then
  NODE_OPTIONS+=("--log_suspect")
else
  NODE_OPTIONS+=("--no-log_suspect")
fi

if [ -n "$PROF" ] && [ $PROF == "y" ]; then
  NODE_OPTIONS+=("--prof")
else
  NODE_OPTIONS+=("--no-prof")
fi

if [ -n "$PROF_CPP" ] && [ $PROF_CPP == "y" ]; then
  NODE_OPTIONS+=("--prof_cpp")
else
  NODE_OPTIONS+=("--no-prof_cpp")
fi

if [ -n "$PROF_BROWSER_MODE" ] && [ $PROF_BROWSER_MODE == "y" ]; then
  NODE_OPTIONS+=("--prof_browser_mode")
else
  NODE_OPTIONS+=("--no-prof_browser_mode")
fi

if [ -n "$LOGFILE_PER_ISOLATE" ] && [ $LOGFILE_PER_ISOLATE == "y" ]; then
  NODE_OPTIONS+=("--logfile_per_isolate")
else
  NODE_OPTIONS+=("--no-logfile_per_isolate")
fi

if [ -n "$LL_PROF" ] && [ $LL_PROF == "y" ]; then
  NODE_OPTIONS+=("--ll_prof")
else
  NODE_OPTIONS+=("--no-ll_prof")
fi

if [ -n "$PERF_BASIC_PROF" ] && [ $PERF_BASIC_PROF == "y" ]; then
  NODE_OPTIONS+=("--perf_basic_prof")
else
  NODE_OPTIONS+=("--no-perf_basic_prof")
fi

if [ -n "$PERF_BASIC_PROF_ONLY_FUNCTIONS" ] && [ $PERF_BASIC_PROF_ONLY_FUNCTIONS == "y" ]; then
  NODE_OPTIONS+=("--perf_basic_prof_only_functions")
else
  NODE_OPTIONS+=("--no-perf_basic_prof_only_functions")
fi

if [ -n "$PERF_PROF" ] && [ $PERF_PROF == "y" ]; then
  NODE_OPTIONS+=("--perf_prof")
else
  NODE_OPTIONS+=("--no-perf_prof")
fi

if [ -n "$PERF_PROF_UNWINDING_INFO" ] && [ $PERF_PROF_UNWINDING_INFO == "y" ]; then
  NODE_OPTIONS+=("--perf_prof_unwinding_info")
else
  NODE_OPTIONS+=("--no-perf_prof_unwinding_info")
fi

if [ -n "$LOG_INTERNAL_TIMER_EVENTS" ] && [ $LOG_INTERNAL_TIMER_EVENTS == "y" ]; then
  NODE_OPTIONS+=("--log_internal_timer_events")
else
  NODE_OPTIONS+=("--no-log_internal_timer_events")
fi

if [ -n "$LOG_TIMER_EVENTS" ] && [ $LOG_TIMER_EVENTS == "y" ]; then
  NODE_OPTIONS+=("--log_timer_events")
else
  NODE_OPTIONS+=("--no-log_timer_events")
fi

if [ -n "$LOG_INSTRUCTION_STATS" ] && [ $LOG_INSTRUCTION_STATS == "y" ]; then
  NODE_OPTIONS+=("--log_instruction_stats")
else
  NODE_OPTIONS+=("--no-log_instruction_stats")
fi

if [ -n "$REDIRECT_CODE_TRACES" ] && [ $REDIRECT_CODE_TRACES == "y" ]; then
  NODE_OPTIONS+=("--redirect_code_traces")
else
  NODE_OPTIONS+=("--no-redirect_code_traces")
fi

if [ -n "$PRINT_OPT_SOURCE" ] && [ $PRINT_OPT_SOURCE == "y" ]; then
  NODE_OPTIONS+=("--print_opt_source")
else
  NODE_OPTIONS+=("--no-print_opt_source")
fi

if [ -n "$TRACE_ELEMENTS_TRANSITIONS" ] && [ $TRACE_ELEMENTS_TRANSITIONS == "y" ]; then
  NODE_OPTIONS+=("--trace_elements_transitions")
else
  NODE_OPTIONS+=("--no-trace_elements_transitions")
fi

if [ -n "$TRACE_CREATION_ALLOCATION_SITES" ] && [ $TRACE_CREATION_ALLOCATION_SITES == "y" ]; then
  NODE_OPTIONS+=("--trace_creation_allocation_sites")
else
  NODE_OPTIONS+=("--no-trace_creation_allocation_sites")
fi

if [ -n "$PRINT_CODE_STUBS" ] && [ $PRINT_CODE_STUBS == "y" ]; then
  NODE_OPTIONS+=("--print_code_stubs")
else
  NODE_OPTIONS+=("--no-print_code_stubs")
fi

if [ -n "$TEST_SECONDARY_STUB_CACHE" ] && [ $TEST_SECONDARY_STUB_CACHE == "y" ]; then
  NODE_OPTIONS+=("--test_secondary_stub_cache")
else
  NODE_OPTIONS+=("--no-test_secondary_stub_cache")
fi

if [ -n "$TEST_PRIMARY_STUB_CACHE" ] && [ $TEST_PRIMARY_STUB_CACHE == "y" ]; then
  NODE_OPTIONS+=("--test_primary_stub_cache")
else
  NODE_OPTIONS+=("--no-test_primary_stub_cache")
fi

if [ -n "$TEST_SMALL_MAX_FUNCTION_CONTEXT_STUB_SIZE" ] && [ $TEST_SMALL_MAX_FUNCTION_CONTEXT_STUB_SIZE == "y" ]; then
  NODE_OPTIONS+=("--test_small_max_function_context_stub_size")
else
  NODE_OPTIONS+=("--no-test_small_max_function_context_stub_size")
fi

if [ -n "$PRINT_CODE" ] && [ $PRINT_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_code")
else
  NODE_OPTIONS+=("--no-print_code")
fi

if [ -n "$PRINT_OPT_CODE" ] && [ $PRINT_OPT_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_opt_code")
else
  NODE_OPTIONS+=("--no-print_opt_code")
fi

if [ -n "$PRINT_CODE_VERBOSE" ] && [ $PRINT_CODE_VERBOSE == "y" ]; then
  NODE_OPTIONS+=("--print_code_verbose")
else
  NODE_OPTIONS+=("--no-print_code_verbose")
fi

if [ -n "$PRINT_BUILTIN_CODE" ] && [ $PRINT_BUILTIN_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_builtin_code")
else
  NODE_OPTIONS+=("--no-print_builtin_code")
fi

if [ -n "$SODIUM" ] && [ $SODIUM == "y" ]; then
  NODE_OPTIONS+=("--sodium")
else
  NODE_OPTIONS+=("--no-sodium")
fi

if [ -n "$PRINT_ALL_CODE" ] && [ $PRINT_ALL_CODE == "y" ]; then
  NODE_OPTIONS+=("--print_all_code")
else
  NODE_OPTIONS+=("--no-print_all_code")
fi

if [ -n "$PREDICTABLE" ] && [ $PREDICTABLE == "y" ]; then
  NODE_OPTIONS+=("--predictable")
else
  NODE_OPTIONS+=("--no-predictable")
fi

if [ -n "$SINGLE_THREADED" ] && [ $SINGLE_THREADED == "y" ]; then
  NODE_OPTIONS+=("--single_threaded")
else
  NODE_OPTIONS+=("--no-single_threaded")
fi