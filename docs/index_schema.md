 - PISI
   - DISTRIBUTION
     - SOURCENAME
       - DESCRIPTION (1+)
         - @LANG
	   - VERSION
	   - TYPE
	   - BinaryName
	   - OBSOLETES
		   + PACKAGE (0+)
   * PACKAGE (1+)
	   - NAME
	   - SUMMARY
		   + @xml:lang
	   - Description
		   + @xml:lang
	   - PartOf
	   - LICENSE (1+)
	   - RuntimeDependencies
		   + Dependency (1+)
			   * @releaseFrom (optional)
	   - History
		   + Update (1-10)
			   * @release
			   * Date
			   * Version
			   * Comment
			   * Name
			   * Email
	   - BuildHost
	   - Distribution
	   - DistributionRelease
	   - Architecture
	   - InstalledSize
	   - PackageSize
	   - PackageHash
	   - PackageURI
	   - DeltaPackages
		   + Delta
			   * @releaseFrom
			   * PackageURI
			   * PackageSize
			   * PackageHash
	   - PackageFormat
	   - Source
		   + Name
		   + Packager
			   * Name
			   * Email
   * Component (1+)
	   - Name
	   - LocalName (0+)
		   + @xml:lang
	   - Summary
		   + @xml:lang
	   - Description
		   + @xml:lang
	   - Group
	   - Maintainer
		   + Name
		   + Email
   * Group (1+)
	   - Name
	   - LocalName
		   + @xml:lang
	   - Icon