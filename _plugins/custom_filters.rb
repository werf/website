module Jekyll
  module CustomFilters
    def true_relative_url(path)
        if !path.instance_of? String
            raise "true_relative_url filter failed: unexpected argument #{path}"
        end

        # remove first slash if exist
        page_path_relative = @context.registers[:page]["url"].gsub(%r!^/!, "")
        page_depth = page_path_relative.scan(%r!/!).count
        prefix = ""
        page_depth.times{ prefix = prefix + "../" }
        prefix + path.sub(%r!^/!, "")
    end

   def normalize_version(version)
     return nil if !version
     return version if version.start_with? "v"
     return "v" + version
   end

   def endswith(text, query)
      return text.end_with? query
    end

    def startswith(text, query)
      return text.start_with? query if text
    end

    def sidebar_entry_by_url(url, entries)
      url_parts = url.gsub(%r!^/!, "").split("/")
      sidebar_url = url_parts[2..].join("/")
      return _find_sidebar_entry_by_sidebar_url_in_folder sidebar_url, entries["f"] if entries["f"]
    end

    def _find_sidebar_entry_by_sidebar_url_in_folder(url, folder)
      folder.each do |entry|
        if entry["url"] == url
          return entry
        end

        if entry["f"]
          res = _find_sidebar_entry_by_sidebar_url_in_folder url, entry["f"]
          return res if res
        end

        if entry["sf"]
          res = _find_sidebar_entry_by_sidebar_url_in_folder url, entry["sf"]
          return res if res
        end
      end

      nil
    end

    # get_lang_field_or_raise_error filter returns a field from argument hash
    # returns nil if hash is empty
    # returns hash[page.lang] if hash has the field
    # returns hash["all"] if hash has the field
    # otherwise, raise an error
    def get_lang_field_or_raise_error(hash)
        if !(hash == nil or hash.instance_of? Hash)
            raise "get_lang_field_or_raise_error filter failed: unexpected argument '#{hash}'"
        end

        if hash == nil or hash.length == 0
            return
        end

        lang = @context.registers[:page]["lang"]
        if hash.has_key?(lang)
            return hash[lang]
        elsif hash.has_key?("all")
            return hash["all"]
        else
            raise "get_lang_field_or_raise_error filter failed: the argument '#{hash}' does not have '#{lang}' or 'all' field"
        end
    end
  end
end

Liquid::Template.register_filter(Jekyll::CustomFilters)
