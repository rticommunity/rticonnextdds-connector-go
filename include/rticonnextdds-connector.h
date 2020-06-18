/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/
#include "lua_binding/lua_binding_ddsConnector.h"
typedef struct RTIDDSConnector RTI_Connector;

int RTI_Connector_get_sample_count(
	void *self,
	const char *entity_name,
	double *value);

int RTI_Connector_get_boolean_from_infos(
	void *self,
	int *return_value,
	const char *entity_name,
	int index,
	const char *name);

int RTI_Connector_set_json_instance(
	void *self,
	const char *entity_name,
	const char *json);

int RTI_Connector_set_boolean_into_samples(
	void *self,
	const char *entity_name,
	const char *name,
	int value);

int RTI_Connector_set_number_into_samples(
	void *self,
	const char *entity_name,
	const char *name,
	double value);

int RTI_Connector_set_string_into_samples(
	void *self,
	const char *entity_name,
	const char *name,
	const char *value);

int RTI_Connector_get_json_from_infos(
	void *self,
	const char *entity_name,
	int index,
	const char *name,
	char **value);

int RTI_Connector_get_json_sample(
	void *self,
	const char *entity_name,
	int index,
	char **json_str);

int RTI_Connector_get_json_member(
	void *self,
	const char *entity_name,
	int index,
	const char *member_name,
	char **json_str);

int RTI_Connector_clear(
	void *self,
	const char *entity_name);

int RTI_Connector_read(
	void *self,
	const char *entity_name);

int RTI_Connector_take(
	void *self,
	const char *entity_name);

int RTI_Connector_write(
	void *self,
	const char *entity_name,
	const char *params_json);

struct RTI_Connector_Options {
        /* boolean */ int enable_on_data_event;
	/* boolean */ int one_based_sequence_indexing;
};


#define RTI_Connector_Options_INITIALIZER { \
        1, /* enable_on_data_event */ \
	1  /* one_based_sequence_indexing */ \
}

RTI_Connector *RTI_Connector_new(
	const char *config_name,
	const char *config_file,
	const struct RTI_Connector_Options *options);

void RTI_Connector_delete(RTI_Connector *self);

int RTI_Connector_get_number_from_sample(
	void *self,
	double *return_value,
	const char *entity_name,
	int index,
	const char *name);

int RTI_Connector_get_boolean_from_sample(
	void *self,
	int *return_value,
	const char *entity_name,
	int index,
	const char *name);

int RTI_Connector_get_string_from_sample(
	void *self,
	char **return_value,
	const char *entity_name,
	int index,
	const char *name);
/*
int RTI_Connector_get_any_from_sample(
	void *self,
	double *double_value_out,
	RTIBool *bool_value_out,
	char **string_value_out,
	RTI_Connector_AnyValueKind *selected_out,
	const char *entity_name,
	int index,
	const char *name);

int RTI_Connector_get_any_from_info(
	void *self,
	double *double_value_out,
	RTIBool *bool_value_out,
	char **string_value_out,
	RTI_Connector_AnyValueKind *selected_out,
	const char *entity_name,
	int index,
	const char *name);
 */

int RTI_Connector_clear_member(
	void *self,
	const char *entity_name,
	const char *name);

void* RTI_Connector_get_datareader( //DDS_DynamicDataReader
	void *self,
	const char *entity_name);

void* RTI_Connector_get_datawriter( //DDS_DynamicDataWriter
	void *self,
	const char *entity_name);

const void* RTI_Connector_get_native_sample( //DDS_DynamicData
	void *self,
	const char *entity_name,
	int index);

int RTI_Connector_wait_for_data(void *self, int timeout);

int RTI_Connector_wait_for_data_on_reader(
	void *self,
	int ms_timeout);

int RTI_Connector_wait_for_acknowledgments(
	void *writer,
	int timeout);

int RTI_Connector_wait_for_matched_publication(
	void *reader,
	int ms_timeout,
	int *current_count_change);

int RTI_Connector_wait_for_matched_subscription(
	void *writer,
	int ms_timeout,
	int *current_count_change);

int RTI_Connector_get_matched_subscriptions(
	void *writer,
	char **json_str);

int RTI_Connector_get_matched_publications(
	void *reader,
	char **json_str);

char * RTI_Connector_get_last_error_message();

int RTI_Connector_get_native_instance(
	void *self,
	const char *entity_name,
	const void **native_pointer); //DDS_DynamicData

void RTI_Connector_free_string(char *str);

int RTI_Connector_set_max_objects_per_thread(
	int value);
