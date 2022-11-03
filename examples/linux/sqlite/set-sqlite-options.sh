#!/bin/bash

set -ex

sqlite_params=""

if [ -n "$AUTOMATIC_INDEX" ]; then
  sqlite_params="${sqlite_params}PRAGMA automatic_index=$AUTOMATIC_INDEX;"
fi

if [ -n "$CASE_SENSITIVE_LIKE" ]; then
  sqlite_params="${sqlite_params}PRAGMA case_sensitive_like=$CASE_SENSITIVE_LIKE;"
fi

if [ -n "$CELL_SIZE_CHECK" ]; then
  sqlite_params="${sqlite_params}PRAGMA cell_size_check=$CELL_SIZE_CHECK;"
fi

if [ -n "$CHECKPOINT_FULLFSYNC" ]; then
  sqlite_params="${sqlite_params}PRAGMA checkpoint_fullfsync=$CHECKPOINT_FULLFSYNC;"
fi

if [ -n "$DEFER_FOREIGN_KEYS" ]; then
  sqlite_params="${sqlite_params}PRAGMA checkpoint_fullfsync=$DEFER_FOREIGN_KEYS;"
fi

if [ -n "$EMPTY_RESULT_CALLBACKS" ]; then
  sqlite_params="${sqlite_params}PRAGMA empty_result_callbacks=$EMPTY_RESULT_CALLBACKS;"
fi

if [ -n "$FOREIGN_KEYS" ]; then
  sqlite_params="${sqlite_params}PRAGMA foreign_keys=$FOREIGN_KEYS;"
fi

if [ -n "$FULL_COLUMN_NAMES" ]; then
  sqlite_params="${sqlite_params}PRAGMA full_column_names=$FULL_COLUMN_NAMES;"
fi

if [ -n "$FULLFSYNC" ]; then
  sqlite_params="${sqlite_params}PRAGMA fullfsync=$FULLFSYNC;"
fi

if [ -n "$PARSER_TRACE" ]; then
  sqlite_params="${sqlite_params}PRAGMA parser_trace=$PARSER_TRACE;"
fi

if [ -n "$QUERY_ONLY" ]; then
  sqlite_params="${sqlite_params}PRAGMA query_only=$QUERY_ONLY;"
fi

if [ -n "$READ_UNCOMMITTED" ]; then
  sqlite_params="${sqlite_params}PRAGMA read_uncommitted=$READ_UNCOMMITTED;"
fi

if [ -n "$RECURSIVE_TRIGGERS" ]; then
  sqlite_params="${sqlite_params}PRAGMA recursive_triggers=$RECURSIVE_TRIGGERS;"
fi

if [ -n "$REVERSE_UNORDERED_SELECTS" ]; then
  sqlite_params="${sqlite_params}PRAGMA reverse_unordered_selects=$REVERSE_UNORDERED_SELECTS;"
fi

if [ -n "$SHORT_COLUMN_NAMES" ]; then
  sqlite_params="${sqlite_params}PRAGMA short_column_names=$SHORT_COLUMN_NAMES;"
fi

if [ -n "$TRUSTED_SCHEMA" ]; then
  sqlite_params="${sqlite_params}PRAGMA trusted_schema=$TRUSTED_SCHEMA;"
fi

if [ -n "$VDBE_ADDOPTRACE" ]; then
  sqlite_params="${sqlite_params}PRAGMA vdbe_addoptrace=$VDBE_ADDOPTRACE;"
fi

if [ -n "$VDBE_DEBUG" ]; then
  sqlite_params="${sqlite_params}PRAGMA vdbe_debug=$VDBE_DEBUG;"
fi

if [ -n "$VDBE_LISTING" ]; then
  sqlite_params="${sqlite_params}PRAGMA vdbe_listing=$VDBE_LISTING;"
fi

if [ -n "$VDBE_TRACE" ]; then
  sqlite_params="${sqlite_params}PRAGMA vdbe_trace=$VDBE_TRACE;"
fi
