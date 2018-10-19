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

void RTIDDSConnector_write(void *connector, char *name, char *param);
void RTIDDSConnector_clear(void *connector, char *name);
void RTIDDSConnector_setNumberIntoSamples(void *connector, char *name, char *field_name, double value);
void RTIDDSConnector_setStringIntoSamples(void *connector, char *name, char *field_name, char *value);
void RTIDDSConnector_setBooleanIntoSamples(void *connector, char *name, char *field_name, int value);
void RTIDDSConnector_setJSONInstance(void *connector, char *name, char *json_str);
void* RTIDDSConnector_getWriter(void *connector, char *name);
void* RTIDDSConnector_getReader(void *connector, char *name);
void RTIDDSConnector_read(void *connector, char *name);
void RTIDDSConnector_take(void *connector, char *name);
double RTIDDSConnector_getSamplesLength(void *connector, char *name);
double RTIDDSConnector_getInfosLength(void *connector, char *name);
int RTIDDSConnector_getBooleanFromInfos(void *connector, char *name, int index, char *field_name);
void* RTIDDSConnector_getJSONSample(void *connector, char *name, int index);
double RTIDDSConnector_getNumberFromSamples(void *connector, char *name, int index, char *field_name);
int RTIDDSConnector_getBooleanFromSamples(void *connector, char *name, int index, char *field_name);
void* RTIDDSConnector_getStringFromSamples(void *connector, char *name, int index, char *field_name);
char* RTIDDSConnector_freeString(char *json);
int RTIDDSConnector_wait(void *connector, int timeout);

