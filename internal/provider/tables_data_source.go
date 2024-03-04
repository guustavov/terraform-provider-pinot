package provider

import (
	"context"
	"fmt"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tablesDataSource{}
	_ datasource.DataSourceWithConfigure = &tablesDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewTablesDataSource() datasource.DataSource {
	return &tablesDataSource{}
}

// usersDataSource is the data source implementation.
type tablesDataSource struct {
	client *goPinotAPI.PinotAPIClient
}

type tablesDataSourceModel struct {
	Tables []tablesModel `tfsdk:"tables"`
}

type TableSegmentsConfig struct {
	TimeType                  string `tfsdk:"timeType"`
	Replication               string `tfsdk:"replication"`
	TimeColumnName            string `tfsdk:"timeColumnName"`
	SegmentAssignmentStrategy string `tfsdk:"segmentAssignmentStrategy"`
	SegmentPushType           string `tfsdk:"segmentPushType"`
	MinimizeDataMovement      bool   `tfsdk:"minimizeDataMovement"`
}

type TableTenant struct {
	Broker string `tfsdk:"broker"`
	Server string `tfsdk:"server"`
}

type StarTreeIndexConfig struct {
	DimensionsSplitOrder              []string `tfsdk:"dimensionsSplitOrder"`
	SkipStarNodeCreationForDimensions []string `tfsdk:"skipStarNodeCreationForDimensions"`
	FunctionColumnPairs               []string `tfsdk:"functionColumnPairs"`
	MaxLeafRecords                    int      `tfsdk:"maxLeafRecords"`
}

type TierOverwrite struct {
	StarTreeIndexConfigs []StarTreeIndexConfig `tfsdk:"starTreeIndexConfigs"`
}

type TierOverwrites struct {
	HotTier  TierOverwrite `tfsdk:"hotTier"`
	ColdTier TierOverwrite `tfsdk:"coldTier"`
}

type TableIndexConfig struct {
	EnableDefaultStarTree                      bool                  `tfsdk:"enableDefaultStarTree"`
	StarTreeIndexConfigs                       []StarTreeIndexConfig `tfsdk:"starTreeIndexConfigs"`
	TierOverwrites                             TierOverwrites        `tfsdk:"tierOverwrites"`
	EnableDynamicStarTreeCreation              bool                  `tfsdk:"enableDynamicStarTreeCreation"`
	AggregateMetrics                           bool                  `tfsdk:"aggregateMetrics"`
	NullHandlingEnabled                        bool                  `tfsdk:"nullHandlingEnabled"`
	OptimizeDictionary                         bool                  `tfsdk:"optimizeDictionary"`
	OptimizeDictionaryForMetrics               bool                  `tfsdk:"optimizeDictionaryForMetrics"`
	NoDictionarySizeRatioThreshold             float64               `tfsdk:"noDictionarySizeRatioThreshold"`
	RangeIndexVersion                          int                   `tfsdk:"rangeIndexVersion"`
	AutoGeneratedInvertedIndex                 bool                  `tfsdk:"autoGeneratedInvertedIndex"`
	CreateInvertedIndexDuringSegmentGeneration bool                  `tfsdk:"createInvertedIndexDuringSegmentGeneration"`
	LoadMode                                   string                `tfsdk:"loadMode"`
	StreamConfigs                              map[string]string     `tfsdk:"streamConfigs"`
}

type TableMetadata struct {
	CustomConfigs map[string]string `tfsdk:"customConfigs"`
}

type TimestampConfig struct {
	Granulatities []string `tfsdk:"granularities"`
}

type FieldIndexInverted struct {
	Enabled string `tfsdk:"enabled"`
}

type FieldIndexes struct {
	Inverted FieldIndexInverted `tfsdk:"inverted"`
}

type FieldConfig struct {
	Name            string          `tfsdk:"name"`
	EncodingType    string          `tfsdk:"encodingType"`
	IndexType       string          `tfsdk:"indexType"`
	IndexTypes      []string        `tfsdk:"indexTypes"`
	TimestampConfig TimestampConfig `tfsdk:"timestampConfig"`
	Indexes         FieldIndexes    `tfsdk:"indexes"`
}

type TransformConfig struct {
	ColumnName        string `tfsdk:"columnName"`
	TransformFunction string `tfsdk:"transformFunction"`
}

type TableIngestionConfig struct {
	SegmentTimeValueCheckType string            `tfsdk:"segmentTimeValueCheckType"`
	TransformConfigs          []TransformConfig `tfsdk:"transformConfigs.omitempty"`
	ContinueOnError           bool              `tfsdk:"continueOnError"`
	RowTimeValueCheck         bool              `tfsdk:"rowTimeValueCheck"`
}

type TierConfig struct {
	Name                string `tfsdk:"name"`
	SegmentSelectorType string `tfsdk:"segmentSelectorType"`
	SegmentAge          string `tfsdk:"segmentAge"`
	StorageType         string `tfsdk:"storageType"`
	ServerTag           string `tfsdk:"serverTag"`
}

type tablesModel struct {
	TableName        string               `tfsdk:"tableName"`
	TableType        string               `tfsdk:"tableType"`
	SegmentsConfig   TableSegmentsConfig  `tfsdk:"segmentsConfig"`
	Tenants          TableTenant          `tfsdk:"tenants"`
	TableIndexConfig TableIndexConfig     `tfsdk:"tableIndexConfig"`
	Metadata         TableMetadata        `tfsdk:"metadata"`
	FieldConfigList  []FieldConfig        `tfsdk:"fieldConfigList,omitempty"`
	IngestionConfig  TableIngestionConfig `tfsdk:"ingestionConfig,omitempty"`
	TierConfigs      []TierConfig         `tfsdk:"tierConfigs,omitempty"`
	IsDimTable       bool                 `tfsdk:"isDimTable"`
}

// Configure adds the provider configured client to the data source.
func (d *tablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goPinotAPI.PinotAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goPinotAPI.PinotAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *tablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tables"
}

// Schema defines the schema for the data source.
func (d *tablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tables": schema.ListNestedAttribute{
				Description: "The list of tables.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tableName": schema.StringAttribute{
							Description: "The name of the table.",
							Computed:    true,
						},
						"tableType": schema.StringAttribute{
							Description: "The type of the table.",
							Computed:    true,
						},
						"segmentsConfig": schema.ListNestedAttribute{
							Description: "The segments configuration of the table.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"timeType": schema.StringAttribute{
										Description: "The time type of the table.",
										Computed:    true,
									},
									"replication": schema.StringAttribute{
										Description: "The replication of the table.",
										Computed:    true,
									},
									"timeColumnName": schema.StringAttribute{
										Description: "The time column name of the table.",
										Computed:    true,
									},
									"segmentAssignmentStrategy": schema.StringAttribute{
										Description: "The segment assignment strategy of the table.",
										Computed:    true,
									},
									"segmentPushType": schema.StringAttribute{
										Description: "The segment push type of the table.",
										Computed:    true,
									},
									"minimizeDataMovement": schema.BoolAttribute{
										Description: "The minimize data movement of the table.",
										Computed:    true,
									},
								},
							},
						},
						"tenants": schema.ListNestedAttribute{
							Description: "The tenants of the table.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"broker": schema.StringAttribute{
										Description: "The broker of the table.",
										Computed:    true,
									},
									"server": schema.StringAttribute{
										Description: "The server of the table.",
										Computed:    true,
									},
								},
							},
						},
						"tableIndexConfig": schema.ListNestedAttribute{
							Description: "The index configuration of the table.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"enableDefaultStarTree": schema.BoolAttribute{
										Description: "The enable default star tree of the table.",
										Computed:    true,
									},
									"starTreeIndexConfigs": schema.ListNestedAttribute{
										Description: "The list of star tree index configurations.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"dimensionsSplitOrder": schema.ListAttribute{
													Description: "The list of dimensions split order.",
													Computed:    true,
												},
												"skipStarNodeCreationForDimensions": schema.ListAttribute{
													Description: "The list of skip star node creation for dimensions.",
													Computed:    true,
												},
												"functionColumnPairs": schema.ListAttribute{
													Description: "The list of function column pairs.",
													Computed:    true,
												},
												"maxLeafRecords": schema.Int64Attribute{
													Description: "The max leaf records.",
													Computed:    true,
												},
											},
										},
									},
									"tierOverwrites": schema.ListNestedAttribute{
										Description: "The tier overwrites of the table.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"hotTier": schema.ListNestedAttribute{
													Description: "The hot tier of the table.",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"starTreeIndexConfigs": schema.ListNestedAttribute{
																Description: "The list of star tree index configurations.",
																Computed:    true,
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"dimensionsSplitOrder": schema.ListAttribute{
																			Description: "The list of dimensions split order.",
																			Computed:    true,
																		},
																		"skipStarNodeCreationForDimensions": schema.ListAttribute{
																			Description: "The list of skip star node creation for dimensions.",
																			Computed:    true,
																		},
																		"functionColumnPairs": schema.ListAttribute{
																			Description: "The list of function column pairs.",
																			Computed:    true,
																		},
																		"maxLeafRecords": schema.Int64Attribute{
																			Description: "The max leaf records.",
																			Computed:    true,
																		},
																	},
																},
															},
														},
													},
												},
												"coldTier": schema.ListNestedAttribute{
													Description: "The cold tier of the table.",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"starTreeIndexConfigs": schema.ListNestedAttribute{
																Description: "The list of star tree index configurations.",
																Computed:    true,
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"dimensionsSplitOrder": schema.ListAttribute{
																			Description: "The list of dimensions split order.",
																			Computed:    true,
																		},
																		"skipStarNodeCreationForDimensions": schema.ListAttribute{
																			Description: "The list of skip star node creation for dimensions.",
																			Computed:    true,
																		},
																		"functionColumnPairs": schema.ListAttribute{
																			Description: "The list of function column pairs.",
																			Computed:    true,
																		},
																		"maxLeafRecords": schema.Int64Attribute{
																			Description: "The max leaf records.",
																			Computed:    true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"enableDynamicStarTreeCreation": schema.BoolAttribute{
										Description: "The enable dynamic star tree creation of the table.",
										Computed:    true,
									},
									"aggregateMetrics": schema.BoolAttribute{
										Description: "The aggregate metrics of the table.",
										Computed:    true,
									},
									"nullHandlingEnabled": schema.BoolAttribute{
										Description: "The null handling enabled of the table.",
										Computed:    true,
									},
									"optimizeDictionary": schema.BoolAttribute{
										Description: "The optimize dictionary of the table.",
										Computed:    true,
									},
									"optimizeDictionaryForMetrics": schema.BoolAttribute{
										Description: "The optimize dictionary for metrics of the table.",
										Computed:    true,
									},
									"noDictionarySizeRatioThreshold": schema.Float64Attribute{
										Description: "The no dictionary size ratio threshold.",
										Computed:    true,
									},
									"rangeIndexVersion": schema.Int64Attribute{
										Description: "The range index version.",
										Computed:    true,
									},
									"autoGeneratedInvertedIndex": schema.BoolAttribute{
										Description: "The auto generated inverted index of the table.",
										Computed:    true,
									},
									"createInvertedIndexDuringSegmentGeneration": schema.BoolAttribute{
										Description: "The create inverted index during segment generation of the table.",
										Computed:    true,
									},
									"loadMode": schema.StringAttribute{
										Description: "The load mode of the table.",
										Computed:    true,
									},
									"streamConfigs": schema.MapAttribute{
										Description: "The stream configs of the table.",
										Computed:    true,
									},
								},
							},
						},
						"metadata": schema.ListNestedAttribute{
							Description: "The metadata of the table.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"customConfigs": schema.MapAttribute{
										Description: "The custom configs of the table.",
										Computed:    true,
									},
								},
							},
						},
						"fieldConfigList": schema.ListNestedAttribute{
							Description: "The list of field configurations.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the field.",
										Computed:    true,
									},
									"encodingType": schema.StringAttribute{
										Description: "The encoding type of the field.",
										Computed:    true,
									},
									"indexType": schema.StringAttribute{
										Description: "The index type of the field.",
										Computed:    true,
									},
									"indexTypes": schema.ListAttribute{
										Description: "The list of index types.",
										Computed:    true,
									},
									"timestampConfig": schema.ListNestedAttribute{
										Description: "The timestamp configuration of the field.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"granularities": schema.ListAttribute{
													Description: "The list of granularities.",
													Computed:    true,
												},
											},
										},
									},
									"indexes": schema.ListNestedAttribute{
										Description: "The indexes of the field.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"inverted": schema.ListNestedAttribute{
													Description: "The inverted index of the field.",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"enabled": schema.StringAttribute{
																Description: "The enabled of the inverted index.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"ingestionConfig": schema.ListNestedAttribute{
							Description: "The ingestion configuration of the table.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"segmentTimeValueCheckType": schema.StringAttribute{
										Description: "The segment time value check type.",
										Computed:    true,
									},
									"transformConfigs": schema.ListNestedAttribute{
										Description: "The list of transform configurations.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"columnName": schema.StringAttribute{
													Description: "The name of the column.",
													Computed:    true,
												},
												"transformFunction": schema.StringAttribute{
													Description: "The transform function.",
													Computed:    true,
												},
											},
										},
									},
									"continueOnError": schema.BoolAttribute{
										Description: "The continue on error.",
										Computed:    true,
									},
									"rowTimeValueCheck": schema.BoolAttribute{
										Description: "The row time value check.",
										Computed:    true,
									},
								},
							},
						},
						"tierConfigs": schema.ListNestedAttribute{
							Description: "The list of tier configurations.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of the tier.",
										Computed:    true,
									},
									"segmentSelectorType": schema.StringAttribute{
										Description: "The segment selector type.",
										Computed:    true,
									},
									"segmentAge": schema.StringAttribute{
										Description: "The segment age.",
										Computed:    true,
									},
									"storageType": schema.StringAttribute{
										Description: "The storage type.",
										Computed:    true,
									},
									"serverTag": schema.StringAttribute{
										Description: "The server tag.",
										Computed:    true,
									},
								},
							},
						},
						"isDimTable": schema.BoolAttribute{
							Description: "The is dim table.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tablesDataSourceModel

	tablesResp, err := d.client.GetTables()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get tables", fmt.Sprintf("Failed to get tables: %s", err))
		return
	}

	for _, tableName := range tablesResp.Tables {
		tableResp, err := d.client.GetTable(tableName)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get table", fmt.Sprintf("Failed to get table: %s", err))
			return
		}

		tableOffline := tablesModel{
			TableName: tableResp.OFFLINE.TableName,
			TableType: tableResp.OFFLINE.TableType,
			SegmentsConfig: TableSegmentsConfig{
				TimeType:                  tableResp.OFFLINE.SegmentsConfig.TimeType,
				Replication:               tableResp.OFFLINE.SegmentsConfig.Replication,
				TimeColumnName:            tableResp.OFFLINE.SegmentsConfig.TimeColumnName,
				SegmentAssignmentStrategy: tableResp.OFFLINE.SegmentsConfig.SegmentAssignmentStrategy,
				SegmentPushType:           tableResp.OFFLINE.SegmentsConfig.SegmentPushType,
				MinimizeDataMovement:      false,
			},
			Tenants: TableTenant{
				Broker: tableResp.OFFLINE.Tenants.Broker,
				Server: tableResp.OFFLINE.Tenants.Server,
			},
			TableIndexConfig: TableIndexConfig{
				EnableDefaultStarTree:                      tableResp.OFFLINE.TableIndexConfig.EnableDefaultStarTree,
				StarTreeIndexConfigs:                       []StarTreeIndexConfig{},
				TierOverwrites:                             TierOverwrites{},
				EnableDynamicStarTreeCreation:              tableResp.OFFLINE.TableIndexConfig.EnableDynamicStarTreeCreation,
				AggregateMetrics:                           tableResp.OFFLINE.TableIndexConfig.AggregateMetrics,
				NullHandlingEnabled:                        tableResp.OFFLINE.TableIndexConfig.NullHandlingEnabled,
				OptimizeDictionary:                         tableResp.OFFLINE.TableIndexConfig.OptimizeDictionary,
				OptimizeDictionaryForMetrics:               tableResp.OFFLINE.TableIndexConfig.OptimizeDictionaryForMetrics,
				NoDictionarySizeRatioThreshold:             tableResp.OFFLINE.TableIndexConfig.NoDictionarySizeRatioThreshold,
				RangeIndexVersion:                          tableResp.OFFLINE.TableIndexConfig.RangeIndexVersion,
				AutoGeneratedInvertedIndex:                 tableResp.OFFLINE.TableIndexConfig.AutoGeneratedInvertedIndex,
				CreateInvertedIndexDuringSegmentGeneration: tableResp.OFFLINE.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration,
				LoadMode:      tableResp.OFFLINE.TableIndexConfig.LoadMode,
				StreamConfigs: tableResp.OFFLINE.TableIndexConfig.StreamConfigs,
			},
			Metadata:        TableMetadata{},
			FieldConfigList: []FieldConfig{},
			IngestionConfig: TableIngestionConfig{},
			TierConfigs:     []TierConfig{},
			IsDimTable:      tableResp.OFFLINE.IsDimTable,
		}

		tableRealtime := tablesModel{
			TableName: tableResp.REALTIME.TableName,
			TableType: tableResp.REALTIME.TableType,
			SegmentsConfig: TableSegmentsConfig{
				TimeType:                  tableResp.REALTIME.SegmentsConfig.TimeType,
				Replication:               tableResp.REALTIME.SegmentsConfig.Replication,
				TimeColumnName:            tableResp.REALTIME.SegmentsConfig.TimeColumnName,
				SegmentAssignmentStrategy: tableResp.REALTIME.SegmentsConfig.SegmentAssignmentStrategy,
				SegmentPushType:           tableResp.REALTIME.SegmentsConfig.SegmentPushType,
				MinimizeDataMovement:      false,
			},
			Tenants: TableTenant{
				Broker: tableResp.REALTIME.Tenants.Broker,
				Server: tableResp.REALTIME.Tenants.Server,
			},
			TableIndexConfig: TableIndexConfig{
				EnableDefaultStarTree:                      tableResp.REALTIME.TableIndexConfig.EnableDefaultStarTree,
				StarTreeIndexConfigs:                       []StarTreeIndexConfig{},
				TierOverwrites:                             TierOverwrites{},
				EnableDynamicStarTreeCreation:              tableResp.REALTIME.TableIndexConfig.EnableDynamicStarTreeCreation,
				AggregateMetrics:                           tableResp.REALTIME.TableIndexConfig.AggregateMetrics,
				NullHandlingEnabled:                        tableResp.REALTIME.TableIndexConfig.NullHandlingEnabled,
				OptimizeDictionary:                         tableResp.REALTIME.TableIndexConfig.OptimizeDictionary,
				OptimizeDictionaryForMetrics:               tableResp.REALTIME.TableIndexConfig.OptimizeDictionaryForMetrics,
				NoDictionarySizeRatioThreshold:             tableResp.REALTIME.TableIndexConfig.NoDictionarySizeRatioThreshold,
				RangeIndexVersion:                          tableResp.REALTIME.TableIndexConfig.RangeIndexVersion,
				AutoGeneratedInvertedIndex:                 tableResp.REALTIME.TableIndexConfig.AutoGeneratedInvertedIndex,
				CreateInvertedIndexDuringSegmentGeneration: tableResp.REALTIME.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration,
				LoadMode:      tableResp.REALTIME.TableIndexConfig.LoadMode,
				StreamConfigs: tableResp.REALTIME.TableIndexConfig.StreamConfigs,
			},
			Metadata:        TableMetadata{},
			FieldConfigList: []FieldConfig{},
			IngestionConfig: TableIngestionConfig{},
			TierConfigs:     []TierConfig{},
			IsDimTable:      tableResp.REALTIME.IsDimTable,
		}

		state.Tables = append(state.Tables, tableOffline, tableRealtime)

	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
