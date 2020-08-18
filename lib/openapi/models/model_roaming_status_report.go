/*
 * Nudm_EE
 *
 * Nudm Event Exposure Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type RoamingStatusReport struct {
	Roaming        bool    `json:"roaming" yaml:"roaming" bson:"roaming" mapstructure:"Roaming"`
	NewServingPlmn *PlmnId `json:"newServingPlmn" yaml:"newServingPlmn" bson:"newServingPlmn" mapstructure:"NewServingPlmn"`
}